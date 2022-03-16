package api

import (
	"net/http"
	"testing"
)

func TestValidateServiceResponse(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "StatusOK",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
		},
		{
			name: "StatusNotFound",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			wantErr: true,
			err:     ErrNotFound,
		},
		{
			name: "StatusTooManyRequests",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusTooManyRequests,
				},
			},
			wantErr: true,
			err:     ErrApiKey,
		},
		{
			name: "SomeOtherStatusCode",
			args: args{
				resp: &http.Response{
					StatusCode: 12345,
				},
			},
			wantErr: true,
			err:     ErrUpstream,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateServiceResponse(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceResponse() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
