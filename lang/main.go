package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"swahili/lang/compiler"
	"swahili/lang/compiler/astformat"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"swahili/lang/server"
)

func main() {
	parseCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	parseCmd.Flags().
		StringP("output", "o", "json", "output format of the tokenizer (json | tree)")

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
		StringP("codegen", "g", "llvm", "code generator to use. llvm or swa")
	compileCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	compileCmd.Flags().
		StringP("output", "o", "start", "location of the compiled executable")
	compileCmd.Flags().
		BoolP("experimental", "e", false, "enable experimental features")

	convertCmd.Flags().
		StringP("source", "s", "", "location of the file containing the source code")
	convertCmd.Flags().
		StringP("target-dialect", "t", "english", "target dialect")
	convertCmd.Flags().
		StringP("out", "o", "", "location of the output file")

	if err := convertCmd.MarkFlagRequired("source"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := convertCmd.MarkFlagRequired("target-dialect"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := compileCmd.MarkFlagRequired("source")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(compileCmd, parseCmd, tokenizeCmd, serverCmd, convertCmd)

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
		codegen, _ := cmd.Flags().GetString("codegen")

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

		switch codegen {
		case "swa":
			req := compiler.ASMCompilerRequest{
				Tree:     &tree,
				Target:   target,
				Dialect:  dialect,
				Filename: source,
			}
			cmp := compiler.NewASMCompiler(&req)

			err := cmp.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "llvm":
			req := compiler.LLVMCompilerRequest{
				Tree:     &tree,
				Target:   target,
				Dialect:  dialect,
				Filename: source,
			}

			err = compiler.Compile(&req)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
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

		for _, v := range tokens {
			fmt.Println(v)
		}
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse the source code",
	Run: func(cmd *cobra.Command, _ []string) {
		source, _ := cmd.Flags().GetString("source")
		output, _ := cmd.Flags().GetString("output")

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

		if output == "tree" {
			dw := astformat.NewTreeDrawer(os.Stdout)
			err = st.Accept(dw)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if output == "json" {

			jf := astformat.NewJsonFormatter()
			err = st.Accept(jf)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			res := jf.Element
			result, err := json.MarshalIndent(res, " ", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(string(result))
		}
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert the source code from one dialect to another",
	Run: func(cmd *cobra.Command, _ []string) {
		sourcePath, _ := cmd.Flags().GetString("source")
		targetDialectName, _ := cmd.Flags().GetString("target-dialect")
		outPath, _ := cmd.Flags().GetString("out")

		// Read source file
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("Error reading source file: %v\n", err)
			os.Exit(1)
		}
		source := string(content)

		// Get target dialect
		targetDialect, ok := lexer.GetDialectByName(targetDialectName)
		if !ok {
			fmt.Printf("Error: Unknown target dialect '%s'\n", targetDialectName)
			os.Exit(1)
		}

		// Initialize lexer with source
		lex, _, err := lexer.NewFastLexer(source)
		if err != nil {
			fmt.Printf("Error initializing lexer: %v\n", err)
			os.Exit(1)
		}

		// Enable whitespace and comments
		lex.ShowWhitespace = true
		lex.ShowComments = true

		tokens, err := lex.GetAllTokens()
		if err != nil {
			fmt.Printf("Error tokenizing source: %v\n", err)
			os.Exit(1)
		}

		// Create a reverse map for the target dialect: TokenKind -> Keyword
		targetKeywords := make(map[lexer.TokenKind]string)
		for keyword, kind := range targetDialect.Reserved() {
			targetKeywords[kind] = keyword
		}

		// 1. Update dialect declaration
		for i := 0; i < len(tokens); i++ {
			if tokens[i].Kind == lexer.DialectDeclaration {
				for j := i + 1; j < len(tokens); j++ {
					k := tokens[j].Kind
					if k == lexer.Whitespace || k == lexer.Comment || k == lexer.Newline {
						continue
					}
					if k == lexer.Colon {
						continue
					}
					if k == lexer.Identifier {
						tokens[j].Raw = targetDialect.Name()
						break
					}
					break
				}
			}
		}

		// 2. Update all translatable tokens (Keywords)
		for i := 0; i < len(tokens); i++ {
			if kw, ok := targetKeywords[tokens[i].Kind]; ok {
				tokens[i].Raw = kw
			}
		}

		// 3. Reconstruct source
		var sb strings.Builder
		for _, token := range tokens {
			sb.WriteString(token.Raw)
		}

		output := sb.String()

		if outPath != "" {
			err := os.WriteFile(outPath, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(output)
		}
	},
}
