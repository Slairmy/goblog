package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPage(t *testing.T) {

	baseURL := "http://127.0.0.1:3000"

	var tests = []struct {
		method   string
		url      string
		excepted int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/1", 200},
		{"GET", "/articles/1/edit", 200},
		{"POST", "/articles/1", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("当前请求: %s", test.url)
		var (
			resp *http.Response
			err  error
		)

		switch test.method {
		case "POST":
			// 构造表单数据
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		case "GET":
			resp, err = http.Get(baseURL + test.url)
		}

		assert.NoError(t, err, "请求"+test.url+"时出错")
		assert.Equal(t, test.excepted, resp.StatusCode, "url"+test.url+"预期出现"+strconv.Itoa(test.excepted)+"实际出现"+strconv.Itoa(resp.StatusCode))

	}
}
