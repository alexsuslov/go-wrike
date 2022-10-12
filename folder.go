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
	"os"
)

const (
	foldersPath      = "/folders"
	folderPath       = "/folders/%s"
	subFoldersPath   = "/folders/%s/folders"
	spaceFoldersPath = "/folders/%s/folders"
)

func GetFolders(ctx context.Context, response *ResponseFolders) error {
	URL := os.Getenv("WRIKE_BASE_URL") + foldersPath
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

// GetFolder by id //https://developers.wrike.com/api/v4/folders-projects/#get-folder
func GetFolder(ctx context.Context, folderId string, response *ResponseFolders) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(folderPath, folderId)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

func GetSubFolders(ctx context.Context, folderId string, response *ResponseFolders) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(subFoldersPath, folderId)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

func GetSpaceFolders(ctx context.Context, spaceId string, response *ResponseFolders) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(spaceFoldersPath, spaceId)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

type ResponseFolders struct {
	Kind string   `json:"kind"`
	Data []Folder `json:"data"`
}

type Folder struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Children []string `json:"children"`
	ChildIds []string `json:"childIds"`
	Scope    string   `json:"scope"`
	Project  Project  `json:"project,omitempty"`
}
