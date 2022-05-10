//Package pget provides the ability to download in parallel.
package pget

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ryojiroakiyama/fileio"
)

//ParallelDownLoadMax is the maximum data size
//that can be downloaded
//by one of the processes running in parallel.
const (
	ParallelDownLoadMax = 1000
)

//Do starts the download from the URL passed as a argument.
//Download process is excuted in parallel.
func Do(url string) error {
	divfiles, err := parallelDownload(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Do: %v", err)
	}
	defer func() {
		for _, d := range divfiles {
			if d != "" {
				fmt.Println("remove:", d)
				os.Remove(d)
			}
		}
	}()
	if err := fileio.BindFiles(divfiles, url[strings.LastIndex(url, "/")+1:]); err != nil {
		return fmt.Errorf("Do: %v", err)
	}
	return nil
}
