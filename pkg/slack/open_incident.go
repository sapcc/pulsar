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

	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/slack/models"
	"github.com/sapcc/pulsar/pkg/util"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &openIncidentCommand{}
	})
}

type openIncidentCommand struct {
	botID string
}

func (o *openIncidentCommand) Init() error {
	cfg, err := config.NewSlackConfigFromEnv()
	if err != nil {
		return err
	}
	o.botID = cfg.BotID
	return nil
}

func (o *openIncidentCommand) IsDisabled() bool {
	//TODO: The incident management is currently in development.
	return true
}

func (o *openIncidentCommand) Describe() string {
	return "Open a new incident"
}

func (o *openIncidentCommand) Keywords() []string {
	return []string{"open incident", "create incident", "new incident"}
}

func (o *openIncidentCommand) RequiredUserRole() auth.UserRole {
	return auth.UserRoles.Base
}

func (o *openIncidentCommand) Run(msg *slack.Msg) (*slack.Msg, error) {
	incident := models.NewIncident(
		util.TrimAnyPrefix(o.Keywords(), msg.Text),
		models.NewUser(models.Reporter, fmt.Sprintf("<@%s>", msg.User)),
		models.NewUser(models.Lead, fmt.Sprintf("<@%s>", o.botID)),
		models.SeverityCritical,
	)

	return incident.ToSlackMessage(), nil
}
