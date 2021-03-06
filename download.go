package pget

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ryojiroakiyama/fileio"

	"golang.org/x/sync/errgroup"
)

func parallelDownload(ctx context.Context, url string) ([]string, error) {
	sumLen, err := checkUrlInfo(url)
	if err != nil {
		return nil, err
	}
	eg, ctx := errgroup.WithContext(ctx)
	nroutine := numOfRoutine(sumLen, 0)
	eachLen := sumLen / int64(nroutine)
	downloadedFiles := make([]string, nroutine)
	for i := 0; i < nroutine; i++ {
		i := i
		err := err
		eg.Go(func() error {
			minRange, maxRange := rangeToDownload(i, nroutine, eachLen, sumLen)
			select {
			case <-ctx.Done(): // Receive cancel and do nothing
			default:
				downloadedFiles[i], err = download(url, minRange, maxRange)
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return downloadedFiles, nil
}

func checkUrlInfo(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("checkUrlInfo: %v", err)
	}
	length := resp.ContentLength
	if length <= 0 {
		return 0, fmt.Errorf("checkUrlInfo: unknown content length")
	}
	// Accept-Ranges value is 'bytes' or 'none' (or omit) define by RFC7233
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, fmt.Errorf("checkUrlInfo: server dosen't support partial requests")
	}
	return length, nil
}

func numOfRoutine(datasize int64, cnt int) int {
	if MaxParallel <= cnt {
		return 0
	}
	if datasize < MinBytesToDownload {
		return 1
	}
	return 1 + numOfRoutine(datasize-MinBytesToDownload, cnt+1)
}

func rangeToDownload(index int, numDiv int, sizeDiv int64, sizeSum int64) (int64, int64) {
	minRange := sizeDiv * int64(index)
	maxRange := sizeDiv * int64(index+1)
	if index == numDiv-1 {
		maxRange += sizeSum - maxRange
	}
	return minRange, maxRange
}

func download(url string, minRange int64, maxRange int64) (string, error) {
	content, err := requestWithRange(url, minRange, maxRange)
	if err != nil {
		return "", err
	}
	defer content.Close()
	return fileio.GenTmpFile(content)
}
