package pget

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var client = &http.Client{}

func rangeValue(start int64, end int64) string {
	return "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
}

// NOTE: rangeの区切りを15個以上にすると503エラー(サーバが過負荷などで処理できない)
func requestWithRange(url string, minRange int64, maxRange int64) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to send request: %v", err)
	}
	req.Header.Add("Range", rangeValue(minRange, maxRange-1))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail to get response: %v", err)
	}
	if resp.StatusCode != http.StatusPartialContent {
		resp.Body.Close()
		return nil, fmt.Errorf("response statuscode : %v", resp.StatusCode)
	}
	return resp.Body, nil
}
