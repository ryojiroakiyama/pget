package pget

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ryojiroakiyama/fileio"

	"golang.org/x/sync/errgroup"
)

func download(ctx context.Context, url string) ([]string, error) {
	eg, ctx := errgroup.WithContext(ctx)
	sumSize, err := dataLengthToDownload(url)
	if err != nil {
		return nil, err
	}
	divNum := numOfRoutine(sumSize)
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

func numOfRoutine(datasize int64) int {
	if datasize < ParallelDownLoadMax {
		return 1
	}
	return 1 + numOfRoutine(datasize/ParallelDownLoadMax)
}

func rangeToDownload(index int, numDiv int, sizeDiv int64, sizeSum int64) (int64, int64) {
	minRange := sizeDiv * int64(index)
	maxRange := sizeDiv * int64(index+1)
	if index == numDiv-1 {
		maxRange += sizeSum - maxRange
	}
	return minRange, maxRange
}

func divDownload(url string, minRange int64, maxRange int64) (string, error) {
	content, err := requestWithRange(url, minRange, maxRange)
	if err != nil {
		return "", err
	}
	defer content.Close()
	return fileio.GenTmpFile(content)
}
