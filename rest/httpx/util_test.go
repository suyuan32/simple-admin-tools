package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	host := "8.8.8.8"
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	r.Header.Set(xForwardedFor, host)
	assert.Equal(t, host, GetRemoteAddr(r))
}

func TestGetRemoteAddrNoHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	assert.True(t, len(GetRemoteAddr(r)) == 0)
}

func TestGetFormValues_TooManyValues(t *testing.T) {
	form := url.Values{}

	// Add more values than the limit
	for i := 0; i < maxFormParamCount+10; i++ {
		form.Add("param", fmt.Sprintf("value%d", i))
	}

	// Create a new request with the form data
	req, err := http.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	assert.NoError(t, err)

	// Set the content type for form data
	req.Header.Set(ContentType, "application/x-www-form-urlencoded")

	_, err = GetFormValues(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too many form values")
}

func TestValidator(t *testing.T) {
	v := NewValidator()
	type User struct {
		Username string `validate:"required,alphanum,max=20"`
		Password string `validate:"required,min=6,max=30"`
	}
	u := User{
		Username: "admin",
		Password: "1",
	}
	result := v.Validate(u, "en")
	assert.Equal(t, "Password must be at least 6 characters in length ", result)

	u = User{
		Username: "admin",
		Password: "123456",
	}
	result = v.Validate(u, "en")
	assert.Equal(t, "", result)
}

func TestParseAcceptLanguage(t *testing.T) {
	data := []struct {
		Str    string
		Target string
	}{
		{
			"zh",
			"zh",
		},
		{
			"zh,en;q=0.9,en-US;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6,la;q=0.5,ja;q=0.4,id;q=0.3,fr;q=0.2",
			"zh",
		},
		{
			"zh-cn,zh;q=0.9",
			"zh",
		},
		{
			"en,zh;q=0.9",
			"en",
		},
	}

	initSupportLanguages()

	for _, v := range data {
		tmp, err := ParseAcceptLanguage(v.Str)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, v.Target, tmp)
	}
}
