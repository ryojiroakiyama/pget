//Package pget provides the ability to download in parallel.
package pget_test

import (
	"os"
	"testing"

	"github.com/ryojiroakiyama/pget"
)

// TODO: Separate Do functions by utilizing interfaces
// TODO: Do test each functions
// TODO: Prepare a test web server
func TestDo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantFile string
	}{
		{
			name: "simple",
			args: args{
				url: "https://github.com/42School/norminette/raw/master/pdf/en.norm.pdf",
			},
			wantErr:  false,
			wantFile: "en.norm.pdf",
		},
		//{
		//	name: "big file",
		//	args: args{
		//		url: "https://ftp.riken.jp/Linux/ubuntu-releases/22.04/ubuntu-22.04-live-server-amd64.iso",
		//		// compare with 'curl -o tmp https://ftp.riken.jp/Linux/ubuntu-releases/22.04/ubuntu-22.04-live-server-amd64.iso'
		//	},
		//	wantErr:  false,
		//	wantFile: "ubuntu-22.04-live-server-amd64.iso",
		//},
		{
			name: "not found",
			args: args{
				url: "no_such",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := pget.Do(tt.args.url); err != nil {
				if !tt.wantErr {
					t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if _, err := os.Stat(tt.wantFile); err != nil {
				t.Errorf("out file dosen't exist, err=%v wantFile=%v", err, tt.wantFile)
				return
			}
			os.Remove(tt.wantFile)
		})
	}
}
