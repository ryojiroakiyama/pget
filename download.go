package pget

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ryojiroakiyama/fileio"

	"golang.org/x/sync/errgroup"
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

func divDownload(url string, minRange int64, maxRange int64) (string, error) {
	content, err := requestWithRange(url, minRange, maxRange)
	if err != nil {
		return "", err
	}
	defer content.Close()
	return fileio.GenTmpFile(content)
}

func download(ctx context.Context, url string) ([]string, error) {
	eg, ctx := errgroup.WithContext(ctx)
	sumSize, err := dataLengthToDownload(url)
	if err != nil {
		return nil, err
	}
	divNum := numOfRangeToDownload(sumSize)
	divSize := sumSize / int64(divNum)
	downloadedFiles := make([]string, divNum)
	for i := 0; i < divNum; i++ {
		i := i
		err := err
		eg.Go(func() error {
			minRange, maxRange := rangeToDownload(i, divNum, divSize, sumSize)
			select {
			case <-ctx.Done(): // Receive cancel and do nothing
			default:
				downloadedFiles[i], err = divDownload(url, minRange, maxRange)
			}
			return err
		})
	}
	// Wait() catch error by Go() if occured, and run the cancellation
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return downloadedFiles, nil
}
