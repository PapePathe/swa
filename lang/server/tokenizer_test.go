package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	web := WebTokenizer{}
	payload := bytes.NewBuffer([]byte(`{"swa": "dialect:french; variable x = 10;"}`))
	request, _ := http.NewRequest(http.MethodPost, "/t", payload)
	response := httptest.NewRecorder()
	expectedResponse := `[{
        "Name": "DIALECT",
        "Value": "dialect"
    },
    {
        "Name": "COLON",
        "Value": ":"
    },
    {
        "Name": "IDENTIFIER",
        "Value": "french"
    },
    {
        "Name": "SEMI_COLON",
        "Value": ";"
    },
    {
        "Name": "LET",
        "Value": "variable"
    },
    {
        "Name": "IDENTIFIER",
        "Value": "x"
    },
    {
        "Name": "ASSIGNMENT",
        "Value": "="
    },
    {
        "Name": "NUMBER",
        "Value": "10"
    },
    {
        "Name": "SEMI_COLON",
        "Value": ";"
    },
    {
        "Name": "EOF",
        "Value": "EOF"
    }]`

	web.ServeHTTP(response, request)

	got := response.Body.String()

	assert.JSONEq(t, got, expectedResponse)
}
