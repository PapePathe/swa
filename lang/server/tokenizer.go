package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"swahili/lang/lexer"
)

type (
	WebTokenizer        struct{}
	WebTokenizerRequest struct {
		Src string `json:"swa"`
	}
)
type WebTokenizerResponse []lexer.Token

func (wt WebTokenizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := wt.expectPostRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

		return
	}

	data, err := wt.expectPayload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	jsonData, err := wt.expectResponse(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if _, err := w.Write(jsonData); err != nil {
		panic(err)
	}
}

func (WebTokenizer) expectPostRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("only post method is allowed")
	}

	return nil
}

func (WebTokenizer) expectPayload(r *http.Request) (*WebParserRequest, error) {
	var data WebParserRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (WebTokenizer) expectResponse(data *WebParserRequest) ([]byte, error) {
	tokens, _, errs := lexer.Tokenize(data.Src)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	jsonData, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
