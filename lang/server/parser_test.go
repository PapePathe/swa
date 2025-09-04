package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEndpoint(t *testing.T) {
	t.Skip()

	web := WebParser{}
	payload := bytes.NewBuffer([]byte(`{"swa": "dialect:french; variable x = 10;"}`))
	request, _ := http.NewRequest(http.MethodPost, "/p", payload)
	response := httptest.NewRecorder()
	expectedResponse := `{
      "ast.BlockStatement": [
        {
          "ast.VarDeclarationStatement": {
            "explicit_type": null,
            "is_constant": false,
            "name": "x",
            "value": {
              "ast.NumberExpression": {
                "value": 10
              }
            }
          }
        }
      ]
    }`

	web.ServeHTTP(response, request)

	got := response.Body.String()

	assert.JSONEq(t, got, expectedResponse)
}
