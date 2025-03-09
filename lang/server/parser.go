package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

type WebParserRequest struct {
	Src string `json:"swa"`
}

type WebParserResponse struct {
	Ast ast.BlockStatement
}
type WebParser struct{}

func (ps WebParser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ps.expectPostRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		return
	}

	data, err := ps.expectPayload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	jsonData, err := ps.expectResponse(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	w.Write(jsonData)
}

func (WebParser) expectPostRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("Only post method is allowed")
	}

	return nil
}

func (WebParser) expectPayload(r *http.Request) (*WebParserRequest, error) {
	var data WebParserRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (WebParser) expectResponse(data *WebParserRequest) ([]byte, error) {
	tokens := lexer.Tokenize(data.Src)
	absTree := parser.Parse(tokens)

	jsonData, err := json.MarshalIndent(absTree, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
