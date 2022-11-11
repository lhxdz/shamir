package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/compute"
	"shamir/pkg/utils/shamir"
)

const (
	// 非流式情况下支持 1MB 数据加密
	stringLimit    = 1 * compute.UnitM
	keyNumberLimit = 1000
)

var (
	fastSplitLen   = compute.GetSecretMaxLen()
	noFastSplitLen = compute.GetSecretMaxLenNoFast()
)

type EncryptCmdConf struct {
	fast   bool
	output string
	t, n   int

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
	cmd.Args = ExactArgs(1)
	// 设置全局flag
	cmd.Flags().BoolVarP(&conf.fast, "fast", "f", true, "Use exist prime to encrypt secret, it will be fast")
	cmd.Flags().StringVarP(&conf.output, "output", "o", "", "Output the keys to path")
	cmd.Flags().IntVarP(&conf.t, "threshold", "t", 0, "The key's threshold, use t keys can decrypt the secret")
	cmd.Flags().IntVarP(&conf.n, "number", "n", 0, "The key's number, this secret will encrypt as n keys")
	cmd.Flags().StringVar(&conf.format, "format", Table, "Output result use [table|yaml|json|csv] "+
		"When use --output, this will not work")

	cmd.RunE = conf.RunE
	return cmd
}

func (e *EncryptCmdConf) RunE(cmd *cobra.Command, args []string) error {
	if len(args[0]) > stringLimit {
		return fmt.Errorf("invalid string, secret length should be less than %d", stringLimit)
	}

	if err := checkTN(e.t, e.n); err != nil {
		return err
	}

	secret := code.EncodeCompoundSecret(args[0], getSplitLen(e.fast))
	keys, prime, err := shamir.HashEncrypt(secret, e.t, e.n, e.fast)
	if err != nil {
		return err
	}

	writer := cmd.OutOrStdout()
	_, err = writer.Write([]byte(fmt.Sprintf("necessary key: %s\n", code.DecodeKeys(prime))))
	if err != nil {
		return err
	}

	header := []string{
		"KEY_X",
		"KEY_Y",
	}
	data := make([][]string, 0, len(keys))
	for _, key := range keys {
		data = append(data, []string{code.DecodeKeys(key.X), code.DecodeKeys(key.Y)})
	}

	return RenderData(e.format, header, data, code.EncodeAbleCompoundKeys(keys), writer)
}

// private

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
