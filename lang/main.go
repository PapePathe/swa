/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"swahili/lang/compiler"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"swahili/lang/server"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "swa",
	Short: "Swahili Programming Environment",
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile the source code to an executable",
	Run: func(cmd *cobra.Command, _ []string) {
		source, _ := cmd.Flags().GetString("source")

		bytes, err := os.ReadFile(source)
		if err != nil {
			os.Stdout.WriteString(err.Error())
			os.Exit(1)
		}
		sourceCode := string(bytes)
		tokens := lexer.Tokenize(sourceCode)
		tree := parser.Parse(tokens)
		target := compiler.BuildTarget{OperatingSystem: "Gnu/Linux", Architecture: "X86-64"}
		compiler.Compile(tree, target)
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
		if err := server.ListenAndServe(); err != nil {
			panic(err)
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
			panic(err)
		}

		sourceCode := string(bytes)
		tokens := lexer.Tokenize(sourceCode)
		st := parser.Parse(tokens)

		result, err := json.MarshalIndent(st, " ", "  ")
		if err != nil {
			panic(err)
		}

		log.Println(string(result))
	},
}

func main() {
	tokenizeCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	tokenizeCmd.Flags().
		StringP("output", "o", "json", "output format of the tokenizer (json | yaml | toml)")

	if err := tokenizeCmd.MarkFlagRequired("source"); err != nil {
		panic(err)
	}

	serverCmd.Flags().
		StringP("server", "l", "", "start a web server")
	compileCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")

	if err := compileCmd.MarkFlagRequired("source"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(compileCmd, tokenizeCmd, serverCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
