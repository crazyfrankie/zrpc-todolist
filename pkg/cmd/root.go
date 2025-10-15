package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/crazyfrankie/zrpc-todolist/pkg/logs"
)

type RootCmd struct {
	Command     *cobra.Command
	processName string
	envPath     string
	serviceName string
}

func NewRootCmd(processName string, serviceName string) *RootCmd {
	rootCmd := &RootCmd{
		processName: processName,
		serviceName: serviceName,
	}
	cmd := &cobra.Command{
		Use:           "Start zRPC-Todolist application",
		Long:          fmt.Sprintf(`Start %s `, processName),
		SilenceUsage:  true,
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd)
		},
	}

	cmd.Flags().StringVarP(&rootCmd.envPath, "env", "e", "", "path of env path")

	rootCmd.Command = cmd
	return rootCmd
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) persistentPreRun(cmd *cobra.Command) error {
	if err := r.initEnv(); err != nil {
		return err
	}
	r.initLog()

	// TODO, other initialize

	return nil
}

func (r *RootCmd) initEnv() error {
	return godotenv.Load(r.envPath)
}

func (r *RootCmd) initLog() {
	logger := logs.NewLogger(os.Stdout)
	logger.WithCaller()
	logger.With("service.name", r.serviceName)
	logs.SetGlobalLogger(logger)
	setLogLevel(logger)
}

func setLogLevel(logger logs.FullLogger) {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch level {
	case "trace":
		logger.SetLevel(logs.LevelTrace)
	case "debug":
		logger.SetLevel(logs.LevelDebug)
	case "info":
		logger.SetLevel(logs.LevelInfo)
	case "notice":
		logger.SetLevel(logs.LevelNotice)
	case "warn":
		logger.SetLevel(logs.LevelWarn)
	case "error":
		logger.SetLevel(logs.LevelError)
	case "fatal":
		logger.SetLevel(logs.LevelFatal)
	default:
		logger.SetLevel(logs.LevelInfo)
	}
}
