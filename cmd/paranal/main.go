package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	meta "github.com/2manyvcos/paranal"
	"github.com/2manyvcos/paranal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   filepath.Base(os.Args[0]),
	Version:               meta.Meta.Version,
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,

	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		return server.Run(cmd.Context())
	},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
