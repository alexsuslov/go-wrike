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
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	CommentsPath       = "/comments"
	taskCommentsPath   = "/tasks/%s/comments"
	folderCommentsPath = "/folders/%s/comments"
)

func GetComments(ctx context.Context, response *ResponseComments) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + CommentsPath
	if err != nil {
		return
	}
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)

}

func getComments(path string) func(ctx context.Context, id string, response *ResponseComments) error {
	return func(ctx context.Context, id string, response *ResponseComments) error {
		URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(path, id)
		body, _, err := Request(ctx, "GET", URL, nil, nil)
		if err != nil {
			return err
		}
		defer body.Close()
		return json.NewDecoder(body).Decode(response)
	}
}

var GetTaskComments = getComments(taskCommentsPath)
var GetFolderComments = getComments(folderCommentsPath)

func CreateTaskComment(ctx context.Context, taskID string, text string, response *ResponseComments) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(taskCommentsPath, taskID)
	if err != nil {
		return
	}
	values := url.Values{}
	values.Add("plainText", "true")
	values.Add("text", text)
	r := io.NopCloser(strings.NewReader(values.Encode()))
	body, _, err := Request(ctx, "POST", URL, r, nil)
	if err != nil {
		return err
	}
	return json.NewDecoder(body).Decode(response)
}

type ResponseComments struct {
	Kind string    `json:"kind"`
	Data []Comment `json:"data"`
}

type Comment struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"authorId"`
	Text        string    `json:"text"`
	UpdatedDate time.Time `json:"updatedDate"`
	CreatedDate time.Time `json:"createdDate"`
	TaskID      string    `json:"taskId,omitempty"`
	FolderID    string    `json:"folderId,omitempty"`
}
