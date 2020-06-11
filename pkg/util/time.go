/*******************************************************************************
*
* Copyright 2019 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package util

import (
	"fmt"
	"time"
)

const timestampFormat = "15:04:05 01.02.2006 UTC"

// HumanizeTimestamp may be used to increase readability of a timestamp.
func HumanizeTimestamp(t time.Time) string {
	return t.Format(timestampFormat)
}

// HumanizeDuration may be used to increase readability of a timestamp.
func HumanizeDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

// StringToTimestamp converts a string to RFC 3339 timestamp.
func StringToTimestamp(theString string) time.Time {
	t, _ := time.Parse(time.RFC3339, theString)
	return t
}

// TimestampToString converts the given time to RFC 3339 format string.
func TimestampToString(ts time.Time) string {
	return ts.Format(time.RFC3339)
}

func TimeStartOfDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 00, 00, now.Location())
}

func TimeEndOfDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 00, 00, now.Location())
}
