package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
