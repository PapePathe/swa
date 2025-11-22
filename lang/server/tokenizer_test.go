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
	payload := bytes.NewBuffer([]byte(`{"swa": "dialecte:français; variable x = 10;"}`))
	request, _ := http.NewRequest(http.MethodPost, "/t", payload)
	response := httptest.NewRecorder()
	expectedResponse := `[{
        "Name": "DIALECT",
        "Value": "dialecte",
        "Line": 0
    },
    {
        "Name": "COLON",
        "Line": 0,
        "Value": ":"
    },
    {
        "Name": "IDENTIFIER",
        "Line": 0,
        "Value": "français"
    },
    {
        "Name": "SEMI_COLON",
        "Line": 0,
        "Value": ";"
    },
    {
        "Name": "LET",
        "Line": 0,
        "Value": "variable"
    },
    {
        "Name": "IDENTIFIER",
        "Line": 0,
        "Value": "x"
    },
    {
        "Name": "ASSIGNMENT",
        "Line": 0,
        "Value": "="
    },
    {
        "Name": "NUMBER",
        "Line": 0,
        "Value": "10"
    },
    {
        "Name": "SEMI_COLON",
        "Line": 0,
        "Value": ";"
    },
    {
        "Name": "EOF",
        "Line": 0,
        "Value": "EOF"
    }]`

	web.ServeHTTP(response, request)

	got := response.Body.String()

	assert.JSONEq(t, got, expectedResponse)
}
