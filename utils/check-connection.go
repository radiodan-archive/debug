package utils

import (
	"io/ioutil"
	"net/http"
)

func CheckConnection() bool {
	// match the request to a known output
	success := "<HTML><HEAD><TITLE>Success</TITLE></HEAD>" +
		"<BODY>Success</BODY></HTML>"

	res, err := http.Get("http://www.apple.com/library/test/success.html")

	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return false
	}

	return string(body) == success
}
