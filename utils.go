package pget

import (
	"fmt"
	"net/http"
	"strconv"
)

func RangeValue(start int64, end int64) string {
	return "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
}

func NumDivideRange(datasize int64) int {
	if datasize < DivDownLoadMax {
		return 1
	}
	return 1 + NumDivideRange(datasize/DivDownLoadMax)
}

func DataLength(url string) (int64, error) {
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

func downloadRange(index int, numDiv int, sizeDiv int64, sizeSum int64) (int64, int64) {
	minRange := sizeDiv * int64(index)
	maxRange := sizeDiv * int64(index+1)
	if index == numDiv-1 {
		maxRange += sizeSum - maxRange
	}
	return minRange, maxRange
}
