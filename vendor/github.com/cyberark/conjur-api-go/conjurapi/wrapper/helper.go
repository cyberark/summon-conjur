package wrapper

import (
	"net/http"
	"io/ioutil"
)

func ByteResponseTransformer(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseText, err
}
