package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hashCommand)
}

var hashCommand = &cobra.Command{
	Use:                   "hash",
	Short:                 "Hash passwords from STDIN line by line",
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    true,

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		hash()
	},
}

func hash() {
	scanner := bufio.NewScanner(os.Stdin)

	noInput := true
	defer func() {
		if noInput {
			fmt.Fprintln(os.Stderr, "Generating hash failed - no input")
			os.Exit(1)
		}
	}()

	for scanner.Scan() {
		credentials := scanner.Text()

		if strings.TrimSpace(credentials) == "" {
			continue
		}
		noInput = false

		hash, err := crypto.Hash(credentials)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while generating hash - %s\n", err)
			os.Exit(1)
		}

		fmt.Println(hash)
	}
}
