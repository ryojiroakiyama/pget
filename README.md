# pget
平行処理によるダウンロード機能を備えたパッケージです。

Package pget provides the ability to download in parallel.

CONSTANTS

const (
        DivDownLoadMax = 1000
)
    DivDownLoadMax is the maximum data size that can be downloaded by one of the
    processes running in parallel.


FUNCTIONS

func DataLength(url string) (int64, error)
func Do(url string) error
    Do starts the download from the URL passed as a argument. Download process
    is excuted in parallel.

func NumDivideRange(datasize int64) int
func RangeValue(start int64, end int64) string
