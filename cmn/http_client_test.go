package cmn

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestSendHttpRequest(t *testing.T) {
	// Mock server that can be configured for different test cases
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			if r.Method == "POST" {
				// Check body for post request

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("Error reading request body: %v", err)
					return
				}

				fmt.Printf("Body: %s\n", string(body))
				if string(body) != `{"key":"val"}` {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			// Check for custom header
			if r.Header.Get("X-Custom-Header") == "custom-value" {
				w.Header().Set("X-Test-Response", "header-received")
			}
			w.WriteHeader(http.StatusOK)
			n, err := fmt.Fprint(w, `{"status":"ok"}`)
			if err != nil {
				panic(err)
			}
			t.Logf("Sent response: %d\n", n)
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"internal server error"}`)
		case "/timeout":
			time.Sleep(100 * time.Millisecond) // Sleep longer than the test timeout
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"status":"delayed"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	type args struct {
		method  string
		url     string
		body    []byte
		headers map[string]string
		timeout time.Duration
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   bool
		errStatus int // Expected status code in AppError
	}{
		{
			name: "Success GET",
			args: args{
				method:  "GET",
				url:     server.URL + "/success",
				timeout: 2 * time.Second,
			},
			want:    []byte(`{"status":"ok"}`),
			wantErr: false,
		},
		{
			name: "Success POST",
			args: args{
				method:  "POST",
				url:     server.URL + "/success",
				body:    []byte(`{"key":"val"}`),
				timeout: 2 * time.Second,
			},
			want:    []byte(`{"status":"ok"}`),
			wantErr: false,
		},
		{
			name: "Server Error",
			args: args{
				method:  "GET",
				url:     server.URL + "/error",
				timeout: 2 * time.Second,
			},
			want:      nil,
			wantErr:   true,
			errStatus: http.StatusInternalServerError,
		},
		{
			name: "Request Timeout",
			args: args{
				method:  "GET",
				url:     server.URL + "/timeout",
				timeout: 50 * time.Millisecond, // Shorter than server sleep
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Custom Header",
			args: args{
				method: "GET",
				url:    server.URL + "/success",
				headers: map[string]string{
					"X-Custom-Header": "custom-value",
				},
				timeout: 2 * time.Second,
			},
			want:    []byte(`{"status":"ok"}`),
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				method:  "GET",
				url:     "invalid-url",
				timeout: 2 * time.Second,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unsupported Scheme",
			args: args{
				method:  "GET",
				url:     "ftp://example.com",
				timeout: 2 * time.Second,
			},
			want:    nil,
			wantErr: true,
		},
	}
	err := ViperInit(".conf_linux.json")
	if err != nil {
		panic(err)
	}
	err = LoggerInit()
	if err != nil {
		panic(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendHttpRequest(tt.args.method, tt.args.url, tt.args.body, tt.args.headers, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendHttpRequest() error = %v, wantErr %v", err.Error(), tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errStatus != 0 {
					var appErr *AppError
					if errors.As(err, &appErr) {
						if appErr.StatusCode != tt.errStatus {
							t.Errorf("SendHttpRequest() error status = %d, wantStatus %d", appErr.StatusCode, tt.errStatus)
						}
					}
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendHttpRequest() = %s, want %s", got, tt.want)
			}
		})
	}
}
