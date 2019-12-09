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

package models

import "fmt"

const (
	emojiReporter = ":man-raising-hand:"
	emojiLead     = ":male-firefighter:"
)

const (
	Reporter = iota
	Lead
)

type (
	UserRole int

	User struct {
		role        UserRole
		displayName string
	}
)

func (r UserRole) String() string {
	return [...]string{"Reporter", "Lead"}[r]
}

func NewUser(role UserRole, displayName string) *User {
	return &User{
		role:        role,
		displayName: displayName,
	}
}

func (u *User) String() string {
	var emoji = emojiReporter
	switch u.role {
	case Reporter:
		emoji = emojiReporter
	case Lead:
		emoji = emojiLead
	}

	return fmt.Sprintf("%s %s: %s", emoji, u.role.String(), u.displayName)
}
