package json

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetResponse(t *testing.T) {
	client := NewClient(resty.New())
	data := client.R().Get("https://jsonplaceholder.typicode.com/todos/1").D()
	assert.Equal(t, 1, data.GetInt("id"))
	assert.Equal(t, false, data.GetBool("completed"))
	assert.Equal(t, "delectus aut autem", data.GetString("title"))
}

func TestPostResponse(t *testing.T) {
	client := NewClient(resty.New())
	data := client.R().J(`
	{
		"title": "foo",
		"body": "bar",
		"userId": 1
	}
	`).Post("https://jsonplaceholder.typicode.com/posts").D()

	assert.Equal(t, 101, data.GetInt("id"))
	assert.Equal(t, 1, data.GetInt("userId"))
	assert.Equal(t, "foo", data.GetString("title"))
}
