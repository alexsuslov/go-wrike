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
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

const (
	jobrolesPath = "/jobroles"
	jobrolePath  = "/jobroles/%s"
)

func GetJobRoles(ctx context.Context, values url.Values, response *ResponseJobRoles) (err error) {
	u, err := url.Parse(os.Getenv("WRIKE_BASE_URL") + jobrolesPath)
	if err != nil {
		return
	}
	u.RawQuery = values.Encode()
	body, _, err := Request(ctx, "GET", u.String(), nil, nil)
	if err != nil {
		return
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

func GetJobRole(ctx context.Context, id string, response *ResponseJobRoles) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(jobrolePath, id)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)

}

type ResponseJobRoles struct {
	Kind string    `json:"kind"`
	Data []JobRole `json:"data"`
}

type JobRole struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ShortTitle string `json:"shortTitle"`
	AvatarURL  string `json:"avatarUrl"`
	IsDeleted  bool   `json:"isDeleted"`
}
