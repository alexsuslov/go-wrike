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
	"time"
)

const (
	taskAttachments        = "/tasks/%s/attachments"
	folderAttachments      = "/tasks/%s/attachments"
	attachmentDownloadPath = "/attachments/%s/download"
)

func DownloadAttachment(ctx context.Context, attachID string) (body io.ReadCloser, err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(attachmentDownloadPath, attachID)
	body, _, err = Request(ctx, "GET", URL, nil, nil)
	return
}

func CreateTaskAttachment(ctx context.Context, taskID string, filename string, body io.ReadCloser) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(taskAttachments, taskID)
	head := map[string]string{
		"Content-Type":     "application/octet-stream",
		"X-Requested-With": "XMLHttpRequest",
		"X-File-Name":      url.QueryEscape(filename),
	}
	_, _, err := Request(ctx, "GET", URL, body, head)
	return err
}

func GetTaskAttachments(ctx context.Context, taskID string, response *ResponseAttachments) (err error) {
	u, err := url.Parse(os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(taskAttachments, taskID))
	if err != nil {
		return
	}
	body, _, err := Request(ctx, "GET", u.String(), nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(body)
}

func GetFolderAttachments(ctx context.Context, folderID string, response *ResponseAttachments) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(folderAttachments, folderID)

	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(body)
}

type ResponseAttachments struct {
	Kind string       `json:"kind"`
	Data []Attachment `json:"data"`
}

type Attachment struct {
	ID                  string    `json:"id"`
	AuthorID            string    `json:"authorId"`
	Name                string    `json:"name"`
	CreatedDate         time.Time `json:"createdDate"`
	Version             int       `json:"version"`
	Type                string    `json:"type"`
	ContentType         string    `json:"contentType"`
	Size                int       `json:"size"`
	TaskID              string    `json:"taskId,omitempty"`
	CurrentAttachmentID string    `json:"currentAttachmentId,omitempty"`
	FolderID            string    `json:"folderId,omitempty"`
}
