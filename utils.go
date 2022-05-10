package pget

import (
	"fmt"
	"net/http"
)

func numOfRangeToDownload(datasize int64) int {
	if datasize < DivDownLoadMax {
		return 1
	}
	return 1 + numOfRangeToDownload(datasize/DivDownLoadMax)
}

func dataLengthToDownload(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("DataLength: %v", err)
	}
	length := resp.ContentLength
	if length <= 0 {
		return 0, fmt.Errorf("DataLength: unknown content length")
	}
	return length, nil
}

func rangeToDownload(index int, numDiv int, sizeDiv int64, sizeSum int64) (int64, int64) {
	minRange := sizeDiv * int64(index)
	maxRange := sizeDiv * int64(index+1)
	if index == numDiv-1 {
		maxRange += sizeSum - maxRange
	}
	return minRange, maxRange
}
