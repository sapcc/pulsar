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

package slack

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/clients"
)

const scheduleName = "Managed Service for CCloud API (Two Day Shifts)"

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &pagerdutyListOnCall{}
	})
}

type pagerdutyListOnCall struct {
	pagerdutyClient *clients.PagerdutyClient
	slackClient     *clients.SlackClient
}

func (l *pagerdutyListOnCall) Init() error {
	pdCli, err := clients.NewPagerdutyClientFromEnv()
	if err != nil {
		return err
	}
	l.pagerdutyClient = pdCli

	sCli, err := clients.NewSlackClientFromEnv()
	if err != nil {
		return err
	}
	l.slackClient = sCli

	return nil
}

func (l *pagerdutyListOnCall) IsDisabled() bool {
	return false
}

func (l *pagerdutyListOnCall) Describe() string {
	return "List on-call persons."
}

func (l *pagerdutyListOnCall) Keywords() []string {
	return []string{"list oncall", "list on-call", "list on call", "who's on call", "who's on-call"}
}

func (l *pagerdutyListOnCall) RequiredUserRole() auth.UserRole {
	return auth.UserRoles.Base
}

func (l *pagerdutyListOnCall) Run(msg *slack.Msg) (*slack.Msg, error) {
	schedule, err := l.pagerdutyClient.GetSchedule(scheduleName)
	if err != nil {
		return nil, err
	}

	onCallUserList, err := l.pagerdutyClient.ListTodaysOnCallUsers(&schedule.ID)
	if err != nil {
		return nil, err
	}

	if len(onCallUserList) == 0 {
		return &slack.Msg{Text: "There's no one on-call right now."}, nil
	}

	users := make([]string, 0)
	for _, u := range onCallUserList {
		if usr, err := l.slackClient.GetUserByEmail(u.Email); err == nil {
			users = append(users, fmt.Sprintf("<@%s>", usr.ID))
		}
	}

	return &slack.Msg{Text: "Currently on call: " + strings.Join(users, ", ")}, nil
}
