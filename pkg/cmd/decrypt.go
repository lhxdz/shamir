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
	"shamir/pkg/utils/path"
	"shamir/pkg/utils/shamir"
)

type DecryptCmdConf struct {
	xKeys []string
	yKeys []string

	necessary         string
	inputPath, output string

	t int
}

func NewDecryptCommand() *cobra.Command {
	cmd := &cobra.Command{}
	conf := &DecryptCmdConf{}
	cmd.Use = "decrypt"
	cmd.Short = "Command line for Shamir decrypt"
	cmd.Long =
		`Command line for Shamir decrypt

You can use it to decrypt n keys which contains (x, y) and one necessary key to secret.
The insertion order of x、y must be the same, and they must be the counts, xKey and yKey will be combined into one key.
`
	cmd.Example = `shamir decrypt -n 123456789 -x 455 -y 455 -x 666 -y 666
shamir decrypt -i ./ -t 2
shamir decrypt -i ./keys/ -t 2 -o ./secret.txt
`
	cmd.Args = NoArgs
	// 设置全局flag
	cmd.Flags().StringVarP(&conf.inputPath, "input-path", "i", "", "The path of keys")
	cmd.Flags().StringVarP(&conf.output, "output", "o", "", "The secret output file")
	cmd.Flags().StringVarP(&conf.necessary, "necessary", "n", "", "The necessary key")
	cmd.Flags().IntVarP(&conf.t, "threshold", "t", 0, "The key's threshold, use t keys to decrypt the secret. "+
		"must use -t when use -p")
	cmd.Flags().StringSliceVarP(&conf.xKeys, "x-key", "x", []string{}, "The key of X")
	cmd.Flags().StringSliceVarP(&conf.yKeys, "y-key", "y", []string{}, "The key of Y")

	cmd.RunE = conf.RunE
	return cmd
}

func (d *DecryptCmdConf) RunE(cmd *cobra.Command, _ []string) error {
	if err := d.check(); err != nil {
		return err
	}

	keyReaders, necessaryReader, taskInputIndicator, err := d.getInput()
	if err != nil {
		return err
	}
	defer taskInputIndicator.Fail()

	output, taskOutputIndicator, err := d.getOutput(cmd)
	if err != nil {
		return err
	}
	defer taskOutputIndicator.Fail()

	keyEncoders := getKeyEncoders(keyReaders)
	nes := code.NewKeyEncoder(necessaryReader)
	secretDecoder := code.NewSecretDecoder(output)

	for {
		keys, necessaryKey, isHash, e := getKeys(keyEncoders, nes)
		if e != nil {
			return e
		}

		e = d.decrypt(keys, necessaryKey, secretDecoder)
		if e != nil {
			return e
		}

		if !isHash {
			continue
		}

		e = secretDecoder.HashCheck()
		if e != nil {
			return e
		}
		break
	}

	// console上的输出换行显示
	if d.output == "" {
		output.Write([]byte("\n"))
	}
	taskInputIndicator.Success()
	taskOutputIndicator.Success()
	return err
}

func (d *DecryptCmdConf) decrypt(keys []code.Key, prime *big.Int, secretWriter *code.SecretDecoder) error {
	secret, e := shamir.Decrypt(keys, prime)
	if e != nil {
		return e
	}

	e = secretWriter.Write(secret)
	if e != nil {
		return e
	}

	return nil
}

func getKeys(keyReaders []*xyKeyEncoder, necessaryReader *code.KeyEncoder) ([]code.Key, *big.Int, bool, error) {
	keys := make([]code.Key, 0, len(keyReaders))
	isHash := false
	for i, reader := range keyReaders {
		key, ok, err := reader.encoder()
		if err != nil {
			return nil, nil, false, err
		}
		if i != 0 && isHash != ok {
			return nil, nil, false, fmt.Errorf("keys not match")
		}

		isHash = ok
		keys = append(keys, key)
	}

	necessaryKey, ok, err := necessaryReader.Read()
	if err != nil {
		return nil, nil, false, err
	}

	if isHash != ok {
		return nil, nil, false, fmt.Errorf("necessary key not match")
	}

	return keys, necessaryKey, isHash, nil
}

func (d *DecryptCmdConf) check() error {
	if d.inputPath != "" {
		if !path.IsExist(d.inputPath) {
			return fmt.Errorf("input path %q not exist", d.inputPath)
		}

		if d.t < shamir.MinThreshold {
			return fmt.Errorf("invalid threshold, please use -t correctly when use input keys by path")
		}
	} else {
		if d.necessary == "" {
			return fmt.Errorf("invalid necessary key, can not be empty")
		}
		if len(d.xKeys) == 0 || len(d.yKeys) == 0 {
			return fmt.Errorf("x keys or y keys can not be zero count")
		}
		if len(d.xKeys) != len(d.yKeys) {
			return fmt.Errorf("xKeys and yKeys must be the same counts")
		}

		for i, xKey := range d.xKeys {
			if xKey == "" || d.yKeys[i] == "" {
				return fmt.Errorf("invalid key, can not have empty key")
			}
		}
	}

	if d.output != "" && path.IsExist(d.output) {
		return fmt.Errorf("output file %q is exist", d.output)
	}

	return nil
}

func (d *DecryptCmdConf) getInput() ([]*keyReadWriter, io.ReadWriter, *TaskIndicator, error) {
	var keys = make([]*keyReadWriter, 0, d.t)
	var necessary io.ReadWriteCloser
	if d.inputPath == "" {
		necessary = NewReadWriteCloser(bytes.NewBufferString(d.necessary))

		for i, xKey := range d.xKeys {
			keys = append(keys, NewKeyReadWriter(NewReadWriteCloser(bytes.NewBufferString(xKey)),
				NewReadWriteCloser(bytes.NewBufferString(d.yKeys[i]))))
		}
		return keys, necessary, NewTaskIndicator(nil, nil), nil
	}

	// 从指定文件夹拿取
	var opened []io.Closer
	d.inputPath = filepath.Clean(d.inputPath)
	keysName, necessaryName, err := path.GetKeysName(d.inputPath)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(keysName) < d.t {
		return nil, nil, nil, fmt.Errorf("invalid input key files, key files can not less than threshold")
	}

	for i := 0; i < d.t; i++ {
		xKeyFileName := filepath.Join(d.inputPath, keysName[i].XKey)
		xKeyFile, err := os.OpenFile(xKeyFileName, os.O_RDONLY, defaultFilePermission)
		if err != nil {
			closeClosers(opened)
			return nil, nil, nil, fmt.Errorf("open x key file %s failed: %w", xKeyFileName, err)
		}
		opened = append(opened, xKeyFile)

		yKeyFileName := filepath.Join(d.inputPath, keysName[i].YKey)
		yKeyFile, err := os.OpenFile(yKeyFileName, os.O_RDONLY, defaultFilePermission)
		if err != nil {
			closeClosers(opened)
			return nil, nil, nil, fmt.Errorf("open y key file %s failed: %w", yKeyFileName, err)
		}
		opened = append(opened, yKeyFile)
		keys = append(keys, NewKeyReadWriter(xKeyFile, yKeyFile))
	}

	necessaryKeyFileName := filepath.Join(d.inputPath, necessaryName)
	necessary, err = os.OpenFile(necessaryKeyFileName, os.O_RDONLY, defaultFilePermission)
	if err != nil {
		closeClosers(opened)
		return nil, nil, nil, fmt.Errorf("open necessary key file %s failed: %w", necessaryKeyFileName, err)
	}
	opened = append(opened, necessary)

	return keys, necessary, NewTaskIndicator(func() { closeClosers(opened) }, func() { closeClosers(opened) }), nil
}

func (d *DecryptCmdConf) getOutput(cmd *cobra.Command) (io.WriteCloser, *TaskIndicator, error) {
	if d.output == "" {
		return NewWriteCloser(cmd.OutOrStdout()), NewTaskIndicator(nil, nil), nil
	}

	d.output = filepath.Clean(d.output)
	if path.IsExist(d.output) {
		return nil, nil, fmt.Errorf("invalid output file path %q, is exist", d.output)
	}

	secret, err := os.OpenFile(d.output, os.O_CREATE|os.O_WRONLY, defaultFilePermission)
	if err != nil {
		return nil, nil, fmt.Errorf("create secret file %q failed: %w", d.output, err)
	}

	return secret, NewTaskIndicator(
		func() { closeClosers([]io.Closer{secret}) },
		func() { rollback([]io.Closer{secret}, []string{d.output}) },
	), nil
}

func getKeyEncoders(keys []*keyReadWriter) []*xyKeyEncoder {
	decoders := make([]*xyKeyEncoder, 0, len(keys))
	for _, key := range keys {
		decoders = append(decoders, key.ToXYKeyEncoder())
	}

	return decoders
}
