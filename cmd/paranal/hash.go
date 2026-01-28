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
			panic(fmt.Errorf("generating hash failed: no input"))
		}
	}()

	for scanner.Scan() {
		credentials := scanner.Text()

		if strings.TrimSpace(credentials) == "" {
			continue
		}
		noInput = false

		hash, err := crypto.Argon2IDHash(credentials)
		if err != nil {
			panic(fmt.Errorf("generating hash failed: %s", err))
		}

		fmt.Println(hash)
	}
}
