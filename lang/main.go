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
