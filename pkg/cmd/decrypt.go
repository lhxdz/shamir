package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"shamir/pkg/utils/code"
	"shamir/pkg/utils/shamir"
)

type DecryptCmdConf struct {
	xKeys []string
	yKeys []string

	necessary string
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
	cmd.Example = "shamir decrypt -n 123456789 -x 455 -y 455 -x 666 -y 666"
	cmd.Args = NoArgs
	// 设置全局flag
	cmd.Flags().StringVarP(&conf.necessary, "necessary", "n", "", "The necessary key")
	cmd.Flags().StringSliceVarP(&conf.xKeys, "x-key", "x", []string{}, "The key of X")
	cmd.Flags().StringSliceVarP(&conf.yKeys, "y-key", "y", []string{}, "The key of Y")

	cmd.RunE = conf.RunE
	return cmd
}

func (d *DecryptCmdConf) RunE(cmd *cobra.Command, _ []string) error {
	if d.necessary == "" {
		return fmt.Errorf("invalid necessary key, can not be empty")
	}
	if len(d.xKeys) == 0 || len(d.yKeys) == 0 {
		return fmt.Errorf("x keys or y keys can not be zero count")
	}
	if len(d.xKeys) != len(d.yKeys) {
		return fmt.Errorf("xKeys and yKeys must be the same counts")
	}

	prime, ok := code.EncodeKeys(d.necessary)
	if !ok {
		return fmt.Errorf("invalid necessary key: %q", d.necessary)
	}
	strKeys := make([]*code.StrKey, 0, len(d.xKeys))
	for i := range d.xKeys {
		strKeys = append(strKeys, &code.StrKey{X: d.xKeys[i], Y: d.yKeys[i]})
	}
	keys, err := code.EncodeStrCompoundKeys(strKeys)
	if err != nil {
		return err
	}

	secret, err := shamir.HashDecrypt(keys, prime)
	if err != nil {
		return err
	}

	_, err = cmd.OutOrStdout().Write([]byte(code.DecodeCompoundSecret(secret) + "\n"))
	return err
}
