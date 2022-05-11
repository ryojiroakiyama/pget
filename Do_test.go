//Package pget provides the ability to download in parallel.
package pget_test

import (
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
		name    string
		args    args
		wantErr bool
	}{
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
			if err := pget.Do(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
