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
        "Name": "",
        "Value": "dialecte",
        "Column": 1,
        "Line": 1,
        "Raw": "dialecte"
    },
    {
        "Name": "",
        "Column": 9,
        "Line": 1,
        "Value": ":",
        "Raw": ":"
    },
    {
        "Name": "",
        "Column": 10,
        "Line": 1,
        "Value": "français",
        "Raw": "français"
    },
    {
        "Name": "",
        "Column": 18,
        "Line": 1,
        "Value": ";",
        "Raw": ";"
    },
    {
        "Name": "",
        "Column": 20,
        "Line": 1,
        "Value": "variable",
        "Raw": "variable"
    },
    {
        "Name": "",
        "Column": 29,
        "Line": 1,
        "Value": "x",
        "Raw": "x"
    },
    {
        "Name": "",
        "Column": 31,
        "Line": 1,
        "Value": "=",
        "Raw": "="
    },
    {
        "Name": "",
        "Column": 33,
        "Line": 1,
        "Value": "10",
        "Raw": "10"
    },
    {
        "Name": "",
        "Column": 35,
        "Line": 1,
        "Value": ";",
        "Raw": ";"
    }]`

	web.ServeHTTP(response, request)

	assert.JSONEq(t, expectedResponse, response.Body.String())
}
