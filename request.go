// (c) 2022 Alex Suslov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package wrike

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Request Request
func Request(ctx context.Context, method string, url string, reader io.ReadCloser,
	head map[string]string) (body io.ReadCloser, header http.Header, err error) {

	InsecureSkipVerify := os.Getenv("InsecureSkipVerify") == "YES"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "bearer "+os.Getenv("WRIKE_API_TOKEN"))
	if head == nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		for k, v := range head {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{Transport: tr}
	r, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("client.Do:%v", err)
		return
	}

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		data := []byte("")
		if r.Body != nil {
			data, _ = io.ReadAll(r.Body)
		}
		err = fmt.Errorf("%v:%v", r.StatusCode, string(data))
		return
	}
	return r.Body, r.Header, err
}
