package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"swahili/lang/compiler"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"swahili/lang/server"
)

func main() {
	parseCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	parseCmd.Flags().
		StringP("output", "o", "json", "output format of the tokenizer (json | yaml | toml)")

	if err := parseCmd.MarkFlagRequired("source"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tokenizeCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	tokenizeCmd.Flags().
		StringP("output", "o", "json", "output format of the tokenizer (json | yaml | toml)")

	if err := tokenizeCmd.MarkFlagRequired("source"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serverCmd.Flags().
		StringP("server", "l", "", "start a web server")

	compileCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	compileCmd.Flags().
		StringP("output", "o", "start", "location of the compiled executable")

	err := compileCmd.MarkFlagRequired("source")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(compileCmd, parseCmd, tokenizeCmd, serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "swa",
	Short: "Swahili Programming Environment",
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile the source code to an executable",
	Run: func(cmd *cobra.Command, _ []string) {
		source, _ := cmd.Flags().GetString("source")
		output, _ := cmd.Flags().GetString("output")

		bytes, err := os.ReadFile(source)
		if err != nil {
			os.Stdout.WriteString(err.Error())
			os.Exit(1)
		}
		sourceCode := string(bytes)
		tokens, dialect := lexer.Tokenize(sourceCode)
		tree, err := parser.Parse(tokens)
		if err != nil {
			os.Stdout.WriteString(err.Error())
			os.Exit(1)
		}
		target := compiler.BuildTarget{
			OperatingSystem: "Gnu/Linux",
			Architecture:    "X86-64",
			Output:          output,
		}
		compiler.Compile(tree, target, dialect)
	},
}

const TIMEOUT = 3

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web service",
	Run: func(_ *cobra.Command, _ []string) {
		mux := http.NewServeMux()
		mux.Handle("/p", server.WebParser{})
		mux.Handle("/t", server.WebTokenizer{})

		server := &http.Server{
			Addr:                         ":2108",
			ReadHeaderTimeout:            TIMEOUT * time.Second,
			Handler:                      mux,
			DisableGeneralOptionsHandler: false,
			WriteTimeout:                 TIMEOUT * time.Second,
		}

		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var tokenizeCmd = &cobra.Command{
	Use:   "tokenize",
	Short: "Tokenize the source code",

	Run: func(cmd *cobra.Command, _ []string) {
		source, _ := cmd.Flags().GetString("source")

		bytes, err := os.ReadFile(source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sourceCode := string(bytes)
		tokens, _ := lexer.Tokenize(sourceCode)

		fmt.Println(tokens)
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse the source code",
	Run: func(cmd *cobra.Command, _ []string) {
		source, _ := cmd.Flags().GetString("source")

		bytes, err := os.ReadFile(source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sourceCode := string(bytes)
		tokens, _ := lexer.Tokenize(sourceCode)
		st, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result, err := json.MarshalIndent(st, " ", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(result))
	},
}
