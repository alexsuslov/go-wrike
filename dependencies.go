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
	taskDependenciesPath = "/tasks/%v/dependencies"
)

func getTaskDependencies(path string) func(ctx context.Context, id string, response *ResponseDependence) error {
	return func(ctx context.Context, id string, response *ResponseDependence) error {
		URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(taskDependenciesPath, id)
		body, _, err := Request(ctx, "GET", URL, nil, nil)
		if err != nil {
			return err
		}
		defer body.Close()
		return json.NewDecoder(body).Decode(response)
	}
}

type ResponseDependence struct {
	Kind string       `json:"kind"`
	Data []Dependence `json:"data"`
}

type Dependence struct {
	ID            string `json:"id"`
	PredecessorID string `json:"predecessorId"`
	SuccessorID   string `json:"successorId"`
	RelationType  string `json:"relationType"`
}
