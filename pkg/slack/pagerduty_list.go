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
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/clients"
	"github.com/sapcc/pulsar/pkg/util"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &pagerdutyList{}
	})
}

type pagerdutyList struct {
	pagerdutyClient *clients.PagerdutyClient
}

func (l *pagerdutyList) Init() error {
	c, err := clients.NewPagerdutyClientFromEnv()
	if err != nil {
		return err
	}
	l.pagerdutyClient = c
	return nil
}

func (l *pagerdutyList) IsDisabled() bool {
	return false
}

func (l *pagerdutyList) Describe() string {
	return "List currently open PagerDuty incidents."
}

func (l *pagerdutyList) Keywords() []string {
	return []string{"list incidents", "incident list"}
}

func (l *pagerdutyList) Run(msg slack.Msg) (slack.Msg, error) {
	f := &clients.IncidentFilter{}
	// If the message contains a cluster name filter incidents accordingly.
	f.ClusterFilterFromText(msg.Text)

	incidentList, err := l.pagerdutyClient.ListIncidents(f)
	if err != nil {
		return slack.Msg{}, err
	}

	if len(incidentList) == 0 {
		response := "No open incidents"
		if f.Clusters != nil {
			response += fmt.Sprintf(" in cluster(s) %s", strings.Join(f.Clusters, ", "))
		}
		response += " :green_heart:"

		return slack.Msg{Text: response}, nil
	}

	data := [][]string{{"Summary", "Started"}}
	for _, inc := range incidentList {
		data = append(data, []string{inc.APIObject.Summary, util.HumanizeTimestamp(util.StringToTimestamp(inc.CreatedAt))})
	}

	return util.ToSlackTable(data), nil
}
