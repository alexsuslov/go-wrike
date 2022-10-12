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
	"time"
)

const (
	invitationsPath = "/invitations"
	invitationPath  = "/invitations/%s"
)

type ResponseInvitations struct {
	Kind string       `json:"kind"`
	Data []Invitation `json:"data"`
}

type Invitation struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"accountId"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	InviterUserID  string    `json:"inviterUserId"`
	InvitationDate time.Time `json:"invitationDate"`
	ResolvedDate   time.Time `json:"resolvedDate"`
	Role           string    `json:"role"`
	External       bool      `json:"external"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
}

func GetInvitations(ctx context.Context, response *ResponseInvitations) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + invitationsPath
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

func GetInvitation(ctx context.Context, contactID string, response *ResponseInvitations) (err error) {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(invitationPath, contactID)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}
