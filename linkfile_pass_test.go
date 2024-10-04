// Note to candidates: Don't change this file.
// When tests in this file pass, you're done!

package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"
)

func TestLinkFile_worksAsIntended(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		want      map[string]*url.URL
		expectErr bool
	}{
		{
			name:     "simple",
			filename: "testdata/valid_file_1.txt",
			want: map[string]*url.URL{
				"example": mustParseURL(t, "https://example.org"),
				"elvia":   mustParseURL(t, "https://elvia.no"),
			},
		},
		{
			name:     "comments",
			filename: "testdata/valid_file_2.txt",
			want: map[string]*url.URL{
				"example": mustParseURL(t, "https://example.org"),
				"elbits":  mustParseURL(t, "https://elbits.no"),
			},
		},
		{
			name:     "invalid",
			filename: "testdata/invalid_file_1.txt",
			want: map[string]*url.URL{
				"notPresentInTheFile": mustParseURL(t, "https://foo"),
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &LinkFile{Filename: test.filename}
			var err error
			for short, want := range test.want {
				got, gErr := f.Long(context.Background(), short)
				if gErr != nil {
					err = errors.Join(err, gErr)
				} else if got.String() != want.String() {
					t.Errorf("Long(%q) = %q, want %q", short, got, want)
				}
			}
			if test.expectErr != (err != nil) {
				t.Errorf("expected error: %t; got %v", test.expectErr, err)
			}
		})
	}

	// Run a small end-to-end test
	t.Run("test server", func(t *testing.T) {
		m := http.NewServeMux()
		m.Handle("GET /{short}", &ShortlinkHandler{&LinkFile{Filename: "testdata/valid_file_1.txt"}})
		s := httptest.NewServer(m)
		t.Cleanup(s.Close)

		ctx, ccl := context.WithTimeout(context.Background(), time.Second*2)
		t.Cleanup(ccl)
		u, err := url.Parse(s.URL)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		u = u.JoinPath("/elvia")
		req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("making HTTP request: %v", err)
		}
		rb, err := httputil.DumpResponse(resp, true)
		if resp.StatusCode >= 400 {
			t.Errorf("reading response: %d", resp.StatusCode)
		}
		t.Log(string(rb))
		if resp.StatusCode >= 400 {
			t.Errorf("bad status code: %d", resp.StatusCode)
		}
		l := resp.Header.Get("Location")
		if l != "https://elvia.no" {
			t.Errorf("Location = %s; want %s", l, "https://elvia.no")
		}
	})
}

func mustParseURL(t testing.TB, s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return u
}
