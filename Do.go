//Package pget provides the ability to download in parallel.
package pget

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ryojiroakiyama/fileio"
)

//ParallelDownLoadMaxLen is the maximum data length
//that can be downloaded
//by one of the processes running in parallel.
const (
	ParallelDownLoadMaxLen = 1000
)

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
