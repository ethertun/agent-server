package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/ethertun/agent-server/server"
	"github.com/ethertun/agent-server/server/endpoints"
	"github.com/spf13/cobra"
)

var (
	// flags
	port      int16
	token     string
	verbosity int

	rootCmd = &cobra.Command{
		Use:   "native-agent",
		Short: "EtherTun local execution agent",
		Long: `An execution agent that runs commands directly in a shell.
                    Useful for testing and debugging.`,
		RunE: run,
	}
)

func init() {
	rootCmd.PersistentFlags().Int16VarP(&port, "port", "p", 1337, "tcp server port to bind")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "unique id/secret used to authenticate with this agent")
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "enables verbose mode (-v, -vv)")

	rootCmd.MarkPersistentFlagRequired("token")
}

func run(cmd *cobra.Command, args []string) error {
	var level slog.Level
	switch v := verbosity; v {
	case 0:
		level = slog.LevelError
	case 1:
		level = slog.LevelInfo
	default:
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	callbacks := server.Callbacks{
		RunTask: run_task,
	}

	server := server.NewServer(port, token, callbacks)
	server.Start()

	return nil
}

func run_task(r *endpoints.RunTaskRequest) (*endpoints.RunTaskOutput, error) {
	slog.Debug("running task", "request", r)

	// split the command into an array
	var stdout, stderr strings.Builder
	args := strings.Split(r.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command failed: %w", err)
	}

	output := &endpoints.RunTaskOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	return output, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to run native-agent", "error", err)
		os.Exit(-1)
	}
}
