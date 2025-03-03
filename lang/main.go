package main

import (
	"fmt"
	"os"
	"swahili/lang/lexer"
	"swahili/lang/parser"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "swa",
	Short: "Swahili Programming Language",
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile the source code to an executable",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Compiling...")
	},
}

var interpretCmd = &cobra.Command{
	Use:   "interpret",
	Short: "Interpret the source code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Interpreting...")
	},
}

var tokenizeCmd = &cobra.Command{
	Use:   "tokenize",
	Short: "Tokenize the source code",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")

		fmt.Println("Tokenizing...")
		bytes, err := os.ReadFile(source)
		if err != nil {
			fmt.Println(fmt.Sprintf("ERROR: %s", err))
			os.Exit(1)
		}

		sourceCode := string(bytes)
		tokens := lexer.Tokenize(sourceCode)
		st := parser.Parse(tokens)

		litter.Dump(st)
	},
}

func main() {
	tokenizeCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	tokenizeCmd.Flags().
		StringP("output", "o", "json", "output format of the tokenizer (json | yaml | toml)")
	tokenizeCmd.MarkFlagRequired("source")

	rootCmd.AddCommand(compileCmd, interpretCmd, tokenizeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
