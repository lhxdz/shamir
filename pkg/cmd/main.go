package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"shamir/pkg/utils/config"
	"shamir/pkg/utils/log"
	"shamir/pkg/version"
)

type MainCmdConf struct {
	console  bool
	logPath  string
	bothLog  bool
	logLevel string
}

func NewMainCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{}
	conf := &MainCmdConf{}
	cmd.Use = "shamir"
	cmd.Short = "Command line for Shamir"
	cmd.Long =
		`Command line for Shamir

"Shamir" be used for (t, n) encrypt. You can use it to encrypt a string or a file.
It will be encrypted as n keys which contains (x, y) and one necessary key.
Any t keys can restore the secret.
For help with any of those, simply call them with --help.`
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}

	// 设置全局flag
	cmd.PersistentFlags().BoolVarP(&conf.console, "console", "c", false, "Print log use console, default (false). "+
		"Can use --both-log and --log print log both to file and console. When use both -c and --log without --both-log, it will only work -c")
	cmd.PersistentFlags().StringVar(&conf.logPath, "log", "", "Print log to log file with path, default empty. "+
		"Can use --both-log and -c print log both to file and console")
	cmd.PersistentFlags().BoolVar(&conf.bothLog, "both-log", false, "Print log both console and file, default (false)")
	cmd.PersistentFlags().StringVarP(&conf.logLevel, "level", "l", "info", "Print log with log level, default (info). "+
		"[debug|info|warn|error|dpanic|panic]")

	cmd.Args = cobra.NoArgs
	cmd.PersistentPreRunE = conf.PreRun

	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Version = version.Version

	// encrypt command
	cmd.AddCommand(NewEncryptCommand())

	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()
	cmd.InitDefaultVersionFlag()

	return cmd, nil
}

func (m *MainCmdConf) PreRun(_ *cobra.Command, _ []string) error {
	err := config.InitConfig(false)
	if err != nil {
		return fmt.Errorf("init config failed: %w", err)
	}

	m.initLog()
	return nil
}

func (m *MainCmdConf) initLog() {
	opts := []log.Option{
		log.WithLogLever(log.Level(m.logLevel)),
		log.WithMaxSize(viper.GetInt(config.LogMaxSize)),
		log.WithMaxBackups(viper.GetInt(config.LogMaxBackups)),
		log.WithMaxAge(viper.GetInt(config.LogMaxAge)),
	}
	if m.logPath == "" {
		m.logPath = viper.GetString(config.LogPathKey)
	}
	if m.console && !m.bothLog {
		m.logPath = ""
	}

	if m.console {
		opts = append(opts, log.WithConsole())
	}
	opts = append(opts, log.WithLogPath(m.logPath))
	log.SetGlobalLogger(log.NewLogger(opts...))
}
