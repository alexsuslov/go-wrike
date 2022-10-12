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
	"net/url"
	"os"
)

const (
	bookingsPath       = "/bookings"
	bookingPath        = "/bookings/%s"
	foldersBookingPath = "/folders/%s/bookings"
)

func GetBookings(ctx context.Context, values url.Values, response *ResponseTasks) (err error) {
	u, err := url.Parse(os.Getenv("WRIKE_BASE_URL") + bookingsPath)
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

type ResponseBookings struct {
	Kind string    `json:"kind"`
	Data []Booking `json:"data"`
}

type Booking struct {
	ID               string                  `json:"id"`
	FolderID         string                  `json:"folderId"`
	ResponsibleID    string                  `json:"responsibleId"`
	BookingDates     BookingDates            `json:"bookingDates"`
	EffortAllocation BookingEffortAllocation `json:"effortAllocation"`
}

type BookingEffortAllocation struct {
	ResponsibleAllocation []ResponsibleAllocation `json:"responsibleAllocation"`
	Mode                  string                  `json:"mode"`
	TotalEffort           int                     `json:"totalEffort"`
}

type ResponsibleAllocation struct {
	UserID          string        `json:"userId"`
	DailyAllocation []interface{} `json:"dailyAllocation"`
}

type BookingDates struct {
	Duration       int    `json:"duration"`
	StartDate      string `json:"startDate"`
	FinishDate     string `json:"finishDate"`
	WorkOnWeekends bool   `json:"workOnWeekends"`
}
