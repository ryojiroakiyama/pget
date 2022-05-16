//Package pget provides the ability to download in parallel.
package pget

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ryojiroakiyama/fileio"
)

//MinBytesToDownload is the minimum data length
//that one of the processes running in parallel downloads.
const MinBytesToDownload = 1 << 7

//MaxParallel is the limit on the number of goroutine.
var MaxParallel = runtime.NumCPU()

//Do starts the download from the URL passed as a argument.
//Download process is excuted in parallel.
func Do(url string) error {
	files, err := parallelDownload(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Do: %v", err)
	}
	defer func() {
		for _, f := range files {
			if f != "" {
				//fmt.Println("remove:", f)
				os.Remove(f)
			}
		}
	}()
	if err := fileio.BindFiles(files, url[strings.LastIndex(url, "/")+1:]); err != nil {
		return fmt.Errorf("Do: %v", err)
	}
	return nil
}
