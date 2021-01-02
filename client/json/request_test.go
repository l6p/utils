package json

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRequest(t *testing.T) {
	client := NewClient(resty.New())
	resp := client.R().Get("https://jsonplaceholder.typicode.com/todos/1")
	assert.Equal(t, "delectus aut autem", resp.D().GetString("title"))
}
