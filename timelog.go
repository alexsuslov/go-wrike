package wrike

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	timelogsPath         = "/timelogs"
	timelogPath          = "/timelogs/%s"
	contactTimelogsPath  = "/contacts/%s/timelogs"
	folderTimelogsPath   = "/folders/%s/timelogs"
	taskTimelogsPath     = "/tasks/%s/timelogs"
	categoryTimelogsPath = "/timelog_categories/%s/timelogs"
)

//https://developers.wrike.com/api/v4/timelogs/

func GetTimeLogs(ctx context.Context, response ResponseTimeLogs) error {
	URL := os.Getenv("WRIKE_BASE_URL") + timelogsPath
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

type ResponseTimeLogs struct {
	Kind string    `json:"kind"`
	Data []TimeLog `json:"data"`
}

type TimeLog struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"taskId"`
	UserID      string    `json:"userId"`
	CategoryID  string    `json:"categoryId,omitempty"`
	Hours       int       `json:"hours"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
	TrackedDate string    `json:"trackedDate"`
	Comment     string    `json:"comment"`
}

func getTimeLogs(path string) func(ctx context.Context, id string, response ResponseTimeLogs) error {
	return func(ctx context.Context, id string, response ResponseTimeLogs) error {
		URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(contactTimelogsPath, id)
		body, _, err := Request(ctx, "GET", URL, nil, nil)
		if err != nil {
			return err
		}
		defer body.Close()
		return json.NewDecoder(body).Decode(response)
	}
}

// GetContactTimeLogs Get all timelog records that were created by the user.
var GetContactTimeLogs = getTimeLogs(contactTimelogsPath)

// GetFolderTimeLogs Get all timelog records for a folder.
var GetFolderTimeLogs = getTimeLogs(folderTimelogsPath)

// GetTaskTimeLogs Get all timelog records for a task.
var GetTaskTimeLogs = getTimeLogs(taskTimelogsPath)

// GetCategoryTimeLogs GetCategoryTimeLogs Get all timelog records with specific timelog category.
var GetCategoryTimeLogs = getTimeLogs(categoryTimelogsPath)

// GetTimeLog Get timelog record by IDs.
var GetTimeLog = getTimeLogs(timelogPath)
