package pget

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func rangeValue(start int64, end int64) string {
	return "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
}

func requestWithRange(url string, minRange int64, maxRange int64) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to send request: %v", err)
	}
	req.Header.Add("Range", rangeValue(minRange, maxRange-1))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail to get response: %v", err)
	}
	return resp.Body, nil
}
