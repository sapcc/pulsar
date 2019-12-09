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
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/clients"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &pagerdutyAcknowledge{}
	})
}

type pagerdutyAcknowledge struct {
	pagerdutyClient *clients.PagerdutyClient
}

func (p *pagerdutyAcknowledge) Init() error {
	client, err := clients.NewPagerdutyClientFromEnv()
	if err != nil {
		return err
	}

	p.pagerdutyClient = client
	return nil
}

func (p *pagerdutyAcknowledge) IsDisabled() bool {
	return true
}

func (p *pagerdutyAcknowledge) Describe() string {
	return "Acknowledge an incident in Pagerduty."
}

func (p *pagerdutyAcknowledge) Keywords() []string {
	return []string{"acknowledge incident", "acknowledge"}
}

func (p *pagerdutyAcknowledge) Run(msg slack.Msg) (slack.Msg, error) {
	return slack.Msg{Text: "I would acknowledge an incident if it would only be implemented."}, nil
}
