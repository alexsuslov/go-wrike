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

const userPath = "/users/%s"

func GetUser(ctx context.Context, userId string, response *ResponseUsers) error {
	URL := os.Getenv("WRIKE_BASE_URL") + fmt.Sprintf(userPath, userId)
	body, _, err := Request(ctx, "GET", URL, nil, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(response)
}

type ResponseUsers struct {
	Kind string `json:"kind"`
	Data []User `json:"data"`
}

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Type      string    `json:"type"`
	Profiles  []Profile `json:"profiles"`
	AvatarURL string    `json:"avatarUrl"`
	Timezone  string    `json:"timezone"`
	Locale    string    `json:"locale"`
	Deleted   bool      `json:"deleted"`
	Me        bool      `json:"me"`
	Phone     string    `json:"phone"`
}
type Profile struct {
	AccountID string `json:"accountId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	External  bool   `json:"external"`
	Admin     bool   `json:"admin"`
	Owner     bool   `json:"owner"`
}
