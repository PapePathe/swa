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
        "Column": 0,
        "Line": 0
    },
    {
        "Name": "COLON",
        "Column": 0,
        "Line": 0,
        "Value": ":"
    },
    {
        "Name": "IDENTIFIER",
        "Column": 0,
        "Line": 0,
        "Value": "français"
    },
    {
        "Name": "SEMI_COLON",
        "Column": 0,
        "Line": 0,
        "Value": ";"
    },
    {
        "Name": "LET",
        "Column": 0,
        "Line": 0,
        "Value": "variable"
    },
    {
        "Name": "IDENTIFIER",
        "Column": 0,
        "Line": 0,
        "Value": "x"
    },
    {
        "Name": "ASSIGNMENT",
        "Column": 0,
        "Line": 0,
        "Value": "="
    },
    {
        "Name": "NUMBER",
        "Column": 0,
        "Line": 0,
        "Value": "10"
    },
    {
        "Name": "SEMI_COLON",
        "Column": 0,
        "Line": 0,
        "Value": ";"
    },
    {
        "Name": "EOF",
        "Column": 0,
        "Line": 0,
        "Value": "EOF"
    }]`

	web.ServeHTTP(response, request)

	got := response.Body.String()

	assert.JSONEq(t, got, expectedResponse)
}
