package cmd

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/compute"
	"shamir/pkg/utils/log"
	"shamir/pkg/utils/path"
	"shamir/pkg/utils/shamir"
)

const (
	// 非流式情况下支持 1MB 数据加密
	stringLimit           = 1 * compute.UnitM
	keyNumberLimit        = 1000
	defaultFilePermission = 0644
)

var (
	// 秘密分隔的大小，应减去额外的前缀开销
	fastSplitLen   = compute.GetSecretMaxLen() - 1
	noFastSplitLen = compute.GetSecretMaxLenNoFast() - 1
)

type EncryptCmdConf struct {
	fast              bool
	outputPath, input string
	t, n              int

	format string
}

func NewEncryptCommand() *cobra.Command {
	cmd := &cobra.Command{}
	conf := &EncryptCmdConf{}
	cmd.Use = "encrypt"
	cmd.Short = "Command line for Shamir encrypt"
	cmd.Long =
		`Command line for Shamir encrypt

You can use it to encrypt a string or a file.
It will be encrypted as n keys which contains (x, y) and one necessary key.
Any t keys can restore the secret.`
	cmd.Example = `shamir encrypt -n 2 -t 2 -o . -i secret.txt
shamir encrypt -n 2 -t 2 -o . < secret.txt
shamir encrypt -n 2 -t 2 "this is a secret.同时支持中文"
`
	// 设置全局flag
	cmd.Flags().BoolVarP(&conf.fast, "fast", "f", true, "Use exist prime to encrypt secret, it will be fast")
	cmd.Flags().StringVarP(&conf.outputPath, "output-path", "o", "", "Output the keys to path")
	cmd.Flags().StringVarP(&conf.input, "input", "i", "", "Read secret from file, if set input file, "+
		"get secret from file first. (must use with -o)")
	cmd.Flags().IntVarP(&conf.t, "threshold", "t", 0, "The key's threshold, use t keys can decrypt the secret")
	cmd.Flags().IntVarP(&conf.n, "number", "n", 0, "The key's number, this secret will encrypt as n keys")
	cmd.Flags().StringVar(&conf.format, "format", Table, "Output result use [table|yaml|json|csv] "+
		"When use --output, this will not work")

	cmd.RunE = conf.RunE
	return cmd
}

func (enc *EncryptCmdConf) RunE(cmd *cobra.Command, args []string) error {
	if err := enc.check(cmd, args); err != nil {
		return err
	}

	input, err := enc.getInput(cmd, args)
	if err != nil {
		return err
	}
	defer closeClosers([]io.Closer{input})
	keys, necessary, taskIndicator, err := enc.getOutput()
	if err != nil {
		return err
	}
	defer taskIndicator.Fail()

	kesDecoders := getKeyDecoders(keys)
	nes := code.NewKeyDecoder(necessary)
	secretReader := code.NewSecretEncoder(input, getSplitLen(enc.fast))
	for {
		subSecret, e := secretReader.Read()
		if e != nil {
			return e
		}
		if subSecret != nil {
			e = enc.encrypt(kesDecoders, nes, subSecret)
			if e != nil {
				return e
			}
		} else {
			e = enc.encrypt(kesDecoders, nes, secretReader.GetHash())
			if e != nil {
				return e
			}
			break
		}
	}
	if enc.outputPath != "" {
		taskIndicator.Success()
		return nil
	}

	writer := cmd.OutOrStdout()
	necessaryData, err := io.ReadAll(necessary)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("necessary key: %s\n", string(necessaryData))))
	if err != nil {
		return err
	}

	header := []string{
		"KEY_X",
		"KEY_Y",
	}
	raw := make([]*code.StrKey, 0, len(keys))
	data := make([][]string, 0, len(keys))
	for _, key := range keys {
		x, y, e := key.toString()
		if e != nil {
			return e
		}
		data = append(data, []string{x, y})
		raw = append(raw, &code.StrKey{X: x, Y: y})
	}

	err = RenderData(enc.format, header, data, raw, writer)
	if err != nil {
		return err
	}

	taskIndicator.Success()
	return nil
}

// private

func (enc *EncryptCmdConf) check(cmd *cobra.Command, args []string) error {
	// terminal 输入秘密
	if enc.input == "" && IsTerminalInput() {
		// 输入参数校验
		if err := ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if len(args[0]) > stringLimit {
			return fmt.Errorf("invalid string, secret length should be less than %dMB, encrypt big secret please use -i", stringLimit/compute.UnitM)
		}
	} else if enc.outputPath == "" {
		return fmt.Errorf("please use -o, when input secret whitout terminal arguments")
	}

	if enc.input != "" && !path.IsExist(enc.input) {
		return fmt.Errorf("invalid input file path %q, not exist", enc.input)
	}

	if err := checkTN(enc.t, enc.n); err != nil {
		return err
	}

	return nil
}

func checkTN(t, n int) error {
	if t < shamir.MinThreshold {
		return fmt.Errorf("invalid threshold %d, should more than %d", t, shamir.MinThreshold)
	}

	if t > n {
		return fmt.Errorf("invalid threshold %d, threshold should less than key number %d", t, n)
	}

	if n > keyNumberLimit {
		return fmt.Errorf("invalid key number %d, key number should less than %d", n, keyNumberLimit)
	}

	return nil
}

func getSplitLen(fast bool) int {
	if fast {
		return fastSplitLen
	}

	return noFastSplitLen
}

func (enc *EncryptCmdConf) getInput(cmd *cobra.Command, args []string) (io.ReadCloser, error) {
	var input io.ReadCloser
	if enc.input == "" {
		// 优先从标准输入拿
		if !IsTerminalInput() {
			return io.NopCloser(cmd.InOrStdin()), nil
		}

		input = io.NopCloser(bytes.NewBufferString(args[0]))
		return input, nil
	}

	enc.input = filepath.Clean(enc.input)
	input, err := os.OpenFile(enc.input, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("open input secret file failed: %w", err)
	}

	return input, nil
}

func (enc *EncryptCmdConf) getOutput() ([]*keyReadWriter, io.ReadWriter, *TaskIndicator, error) {
	var keys = make([]*keyReadWriter, 0, enc.n)
	var necessary io.ReadWriteCloser
	if enc.outputPath == "" {
		necessary = NewReadWriteCloser(bytes.NewBuffer([]byte{}))

		for i := 0; i < enc.n; i++ {
			keys = append(keys, NewKeyReadWriter(NewReadWriteCloser(bytes.NewBuffer([]byte{})),
				NewReadWriteCloser(bytes.NewBuffer([]byte{}))))
		}
		return keys, necessary, NewTaskIndicator(nil, nil), nil
	}

	// 输出到指定文件夹下
	enc.outputPath = filepath.Clean(enc.outputPath)
	err := os.MkdirAll(enc.outputPath, 0750)
	if err != nil {
		return nil, nil, nil, err
	}

	err = path.CheckNoKey(enc.outputPath)
	if err != nil {
		return nil, nil, nil, err
	}

	var opened []io.Closer
	var paths []string
	for i := 0; i < enc.n; i++ {
		xKeyFileName := filepath.Join(enc.outputPath, getXKeyFileName(i))
		xKeyFile, err := os.OpenFile(xKeyFileName, os.O_CREATE|os.O_WRONLY, defaultFilePermission)
		if err != nil {
			rollback(opened, paths)
			return nil, nil, nil, fmt.Errorf("create x key file %s failed: %w", xKeyFileName, err)
		}
		opened = append(opened, xKeyFile)
		paths = append(paths, xKeyFileName)

		yKeyFileName := filepath.Join(enc.outputPath, getYKeyFileName(i))
		yKeyFile, err := os.OpenFile(yKeyFileName, os.O_CREATE|os.O_WRONLY, defaultFilePermission)
		if err != nil {
			rollback(opened, paths)
			return nil, nil, nil, fmt.Errorf("create y key file %s failed: %w", xKeyFileName, err)
		}
		opened = append(opened, yKeyFile)
		paths = append(paths, yKeyFileName)
		keys = append(keys, NewKeyReadWriter(xKeyFile, yKeyFile))
	}

	necessaryKeyFileName := filepath.Join(enc.outputPath, path.NecessaryFileName)
	necessary, err = os.OpenFile(necessaryKeyFileName, os.O_CREATE|os.O_WRONLY, defaultFilePermission)
	if err != nil {
		rollback(opened, paths)
		return nil, nil, nil, fmt.Errorf("create necessary key file %s failed: %w", necessaryKeyFileName, err)
	}
	opened = append(opened, necessary)
	paths = append(paths, necessaryKeyFileName)

	return keys, necessary, NewTaskIndicator(func() { closeClosers(opened) }, func() { rollback(opened, paths) }), nil
}

func getXKeyFileName(id int) string {
	return fmt.Sprintf("%s%d", path.XKeyFilePrefix, id)
}
func getYKeyFileName(id int) string {
	return fmt.Sprintf("%s%d", path.YKeyFilePrefix, id)
}

func rollback(opened []io.Closer, paths []string) {
	closeClosers(opened)
	deleteFiles(paths)
}

func closeClosers(opened []io.Closer) {
	for _, file := range opened {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}
}

func deleteFiles(paths []string) {
	for _, file := range paths {
		err := os.Remove(file)
		if err != nil {
			log.Error(err)
		}
	}
}

func getKeyDecoders(keys []*keyReadWriter) []*xyKeyDecoder {
	decoders := make([]*xyKeyDecoder, 0, len(keys))
	for _, key := range keys {
		decoders = append(decoders, key.ToXYKeyDecoder())
	}

	return decoders
}

func (enc *EncryptCmdConf) encrypt(keys []*xyKeyDecoder, necessary *code.KeyDecoder, secret *big.Int) error {
	subKeys, subPrime, e := shamir.Encrypt(secret, enc.t, enc.n, enc.fast)
	if e != nil {
		return e
	}
	for i, key := range subKeys {
		e = keys[i].decoder(&key)
		if e != nil {
			return e
		}
	}

	e = necessary.Write(subPrime)
	if e != nil {
		return e
	}

	return nil
}
