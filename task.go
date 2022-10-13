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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/url"
	"os"
	"time"
)

const (
	tasksPath            = "/tasks"
	taskPath             = "/tasks/%s"
	folderTasksPath      = "/folders/%s/tasks"
	spaceTasksPath       = "/spaces/%s/tasks"
	folderTaskCreatePath = "/folders/%v/tasks"
)

func GetTasks(ctx context.Context, values url.Values, response *ResponseTasks) (err error) {
	u, err := url.Parse(os.Getenv("WRIKE_BASE_URL") + tasksPath)
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

func getTasks(path string) func(ctx context.Context, id string, response *ResponseTasks) error {
	return func(ctx context.Context, id string, response *ResponseTasks) error {
		URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(path, id)
		body, _, err := Request(ctx, "GET", URL, nil, nil)
		if err != nil {
			return err
		}
		defer body.Close()
		return json.NewDecoder(body).Decode(response)
	}
}

var GetTask = getTasks(taskPath)
var GetFolderTasks = getTasks(folderTasksPath)
var GetSpaceTasks = getTasks(spaceTasksPath)

func FolderCreateTask(ctx context.Context, folderID string, req CreateTaskRQ, resp *ResponseTasks) (err error) {
	urlFormat := os.Getenv("WRIKE_BASE_URL") + folderTaskCreatePath
	url := fmt.Sprintf(urlFormat, folderID)
	Data, err := req.ToRequest(url)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte(Data))
	body, _, err := Request(ctx, "POST", url, io.NopCloser(buf), nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(resp)
}

type ResponseTasks struct {
	Kind string `json:"kind"`
	Data []Task `json:"data"`
}

type Task struct {
	ID               string        `json:"id"`
	AccountID        string        `json:"accountId"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	BriefDescription string        `json:"briefDescription"`
	ParentIds        []string      `json:"parentIds"`
	SuperParentIds   []interface{} `json:"superParentIds"`
	SharedIds        []string      `json:"sharedIds"`
	ResponsibleIds   []string      `json:"responsibleIds"`
	Status           string        `json:"status"`
	Importance       string        `json:"importance"`
	CreatedDate      time.Time     `json:"createdDate"`
	UpdatedDate      time.Time     `json:"updatedDate"`
	Dates            TaskDate      `json:"dates"`
	Scope            string        `json:"scope"`
	AuthorIds        []string      `json:"authorIds"`
	CustomStatusID   string        `json:"customStatusId"`
	HasAttachments   bool          `json:"hasAttachments"`
	Permalink        string        `json:"permalink"`
	Priority         string        `json:"priority"`
	FollowedByMe     bool          `json:"followedByMe"`
	FollowerIds      []string      `json:"followerIds"`
	SuperTaskIds     []interface{} `json:"superTaskIds"`
	SubTaskIds       []string      `json:"subTaskIds"`
	DependencyIds    []string      `json:"dependencyIds"`
	Metadata         Metadatas     `json:"metadata"`
	CustomFields     []Fields      `json:"customFields"`
}

type Project struct {
	AuthorID       string   `json:"authorId"`
	OwnerIds       []string `json:"ownerIds"`
	Status         string   `json:"status"`
	CustomStatusID string   `json:"customStatusId"`
	StartDate      string   `json:"startDate"`
	EndDate        string   `json:"endDate"`
	CreatedDate    string   `json:"createdDate"`
	CompletedDate  string   `json:"completedDate"`
}

type CreateTaskRQ struct {
	Title            string            `json:"title"`
	Description      *string           `json:"description"`
	Status           *string           `json:"status"`
	СustomStatus     *string           `json:"customStatus"`
	Importance       *string           `json:"importance"`
	Dates            *TaskDate         `json:"dates"`
	BillingType      string            `json:"billingType"`
	Shareds          []string          `json:"shareds"`
	Parents          []string          `json:"parents"`
	Responsibles     []string          `json:"responsibles"`
	Followers        []string          `json:"followers"`
	Follow           bool              `json:"follow"`
	PriorityBefore   *string           `json:"priorityBefore"`
	PriorityAfter    *string           `json:"priorityAfter"`
	SuperTasks       *string           `json:"superTasks"`
	Metadata         Metadatas         `json:"metadata"`
	CustomFields     Fields            `json:"customFields"`
	EffortAllocation *EffortAllocation `json:"effortAllocation"`
}

type TaskDate struct {
	Type     string  `json:"type"`
	Duration *int    `json:"duration"`
	Start    *string `json:"start"`
	Due      *string `json:"due"`
}

type Metadatas []Metadata

type Metadata struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

func (Metadatas Metadatas) ToBytes() []byte {
	data, _ := json.Marshal(Metadatas)
	return data
}

type Fields []Field

func (Fields Fields) ToBytes() []byte {
	data, _ := json.Marshal(Fields)
	return data
}

type Field struct {
	ID    string  `json:"id"`
	Value *string `json:"value"`
}

type EffortAllocation struct {
	mode            string
	totalEffort     *int
	allocatedEffort *int
	billingType     *string
	fields          Fields
}

func (t CreateTaskRQ) ToRequest(URL string) (r string, err error) {
	form := url.Values{}
	form.Add("title", html.EscapeString(t.Title))
	if t.Description != nil {
		form.Add("description", *t.Description)
	}
	if t.Status != nil {
		form.Add("status", *t.Status)
	}
	if t.СustomStatus != nil {
		form.Add("customStatus", *t.СustomStatus)
	}

	if t.Metadata != nil {
		form.Add("metadata", string(t.Metadata.ToBytes()))
	}

	if t.CustomFields != nil {
		form.Add("customFields", string(t.CustomFields.ToBytes()))
	}

	if t.Followers != nil && len(t.Followers) > 0 {
		data, _ := json.Marshal(t.Followers)
		form.Add("followers", string(data))
	}

	if t.Responsibles != nil && len(t.Responsibles) > 0 {
		data, _ := json.Marshal(t.Responsibles)
		form.Add("responsibles", string(data))
	}

	return url.QueryUnescape(form.Encode())
}
