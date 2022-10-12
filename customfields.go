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
	customFieldsPath = "/customfields"
	customFieldPath  = "/customfields/%s"
)

func GetCustomFields(ctx context.Context, response *ResponseCustomFields) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + customFieldsPath
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

func GetCustomField(ctx context.Context, customFieldID string, response *ResponseCustomFields) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(customFieldPath, customFieldID)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

type ResponseCustomFields struct {
	Kind string        `json:"kind"`
	Data []Customfield `json:"data"`
}

type Customfield struct {
	ID        string        `json:"id"`
	AccountID string        `json:"accountId"`
	Title     string        `json:"title"`
	Type      string        `json:"type"`
	SharedIds []interface{} `json:"sharedIds"`
	Settings  struct {
		InheritanceType       string `json:"inheritanceType"`
		DecimalPlaces         int    `json:"decimalPlaces"`
		UseThousandsSeparator bool   `json:"useThousandsSeparator"`
		Aggregation           string `json:"aggregation"`
		ReadOnly              bool   `json:"readOnly"`
	} `json:"settings,omitempty"`
	Settings0 struct {
		InheritanceType string `json:"inheritanceType"`
		ReadOnly        bool   `json:"readOnly"`
	} `json:"settings,omitempty"`
	Settings1 struct {
		InheritanceType string `json:"inheritanceType"`
		ReadOnly        bool   `json:"readOnly"`
	} `json:"settings,omitempty"`
}
