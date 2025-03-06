package server

import (
	"encoding/json"
	"net/http"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

type LiveInterpreterRequest struct {
	Source string `json:"swa"`
}
type LiveInterpreterResponse struct {
	Tokens []lexer.Token
	Ast    ast.BlockStatement
}

type LiveInterpreter struct{}

func (liv LiveInterpreter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Source code must be sent though post request", http.StatusMethodNotAllowed)

		return
	}

	var data LiveInterpreterRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)

		return
	}

	tokens := lexer.Tokenize(data.Source)
	absTree := parser.Parse(tokens)

	response := LiveInterpreterResponse{
		//		Tokens: tokens,
		Ast: absTree,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Unable to encode JSON", http.StatusInternalServerError)
	}

	w.Write(jsonData)
}
