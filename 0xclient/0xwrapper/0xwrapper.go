package ox

import (
  "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//returns body of get request to a specified url
func GetRequest(url string) (body []byte, err error) {
	tr := &http.Transport{
		MaxIdleConns:       1,
		IdleConnTimeout:    1 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && err == nil {
		err = errors.New(fmt.Sprintf("Status: %v\r\n", resp.StatusCode))
	}
	resp.Body.Close()
	return
}
