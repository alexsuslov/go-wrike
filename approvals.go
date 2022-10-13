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
	"time"
)

const (
	approvalsPath       = "/approvals"
	approvalPath        = "/approvals/%id"
	folderApprovalsPath = "/folders/%s/approvals"
	taskApprovalsPath   = "/tasks/%s/approvals"
)

func GetApprovals(ctx context.Context, values url.Values, response *ResponseApprovals) (err error) {
	u, err := url.Parse(os.Getenv("WRIKE_BASE_URL") + approvalsPath)
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

func getApprovals(path string) func(ctx context.Context, id string, response *ResponseApprovals) error {
	return func(ctx context.Context, id string, response *ResponseApprovals) error {
		URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(path, id)
		body, _, err := Request(ctx, "GET", URL, nil, nil)
		if err != nil {
			return err
		}
		defer body.Close()
		return json.NewDecoder(body).Decode(response)
	}
}

var GetApproval = getApprovals(approvalPath)
var GetFolderApproval = getApprovals(folderApprovalsPath)
var GetTaskrApproval = getApprovals(taskApprovalsPath)

type ResponseApprovals struct {
	Kind string     `json:"kind"`
	Data []Approval `json:"data"`
}

type Approval struct {
	TaskID              string      `json:"taskId,omitempty"`
	AuthorID            string      `json:"authorId"`
	Title               string      `json:"title"`
	Description         string      `json:"description"`
	UpdatedDate         time.Time   `json:"updatedDate"`
	DueDate             string      `json:"dueDate"`
	Decisions           []Decisions `json:"decisions"`
	ReviewID            string      `json:"reviewId"`
	AttachmentIds       []string    `json:"attachmentIds"`
	Type                string      `json:"type"`
	AutoFinishOnApprove bool        `json:"autoFinishOnApprove"`
	AutoFinishOnReject  bool        `json:"autoFinishOnReject"`
	Finished            bool        `json:"finished"`
	ID                  string      `json:"id"`
	Status              string      `json:"status"`
	FolderID            string      `json:"folderId,omitempty"`
}

type Decisions struct {
	ApproverID  string    `json:"approverId"`
	Comment     string    `json:"comment"`
	Status      string    `json:"status"`
	UpdatedDate time.Time `json:"updatedDate"`
}
