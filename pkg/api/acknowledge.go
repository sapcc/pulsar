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

package api

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/clients"
)

// acknowledge will:
// 1. open a slack thread noting the acknowledger (slack user) and time
// 2. add an emoji to the original slack message with the alert to indicate it's being worked on
// 3. search and acknowledge the corresponding incident in pagerduty to avoid further esacalation (call, etc.)
func (a *API) acknowledge(message slack.InteractionCallback) error {
	// Post the message.
	if _, _, err := a.slackClient.PostMessage(
		message.Channel.ID,
		slack.MsgOptionText(fmt.Sprintf(acknowledgeString, message.User.ID), false),
		slack.MsgOptionTS(message.OriginalMessage.Timestamp),
	); err != nil {
		return err
	}

	// Add reaction emoji to original message.
	if err := a.slackClient.AddReactionToMessage(
		message.Channel.ID,
		message.OriginalMessage.Timestamp,
		emojiFirefighter,
	); err != nil {
		return err
	}

	slackUser, err := a.slackClient.GetUserByID(message.User.ID)
	if err != nil {
		level.Error(a.logger).Log("msg", "cannot find slack user", "err", err.Error())
		return err
	}

	// Find the corresponding pagerduty user.
	user, err := a.pdClient.GetUserByEmail(slackUser.Profile.Email)
	if err != nil {
		level.Info(a.logger).Log("msg", "failed to find pagerduty user. falling back to default user", "err", err.Error())
		user = a.pdClient.GetDefaultUser()
	}

	if len(message.OriginalMessage.Attachments) == 0 ||  message.OriginalMessage.Attachments[0].Text == "" {
		return errors.New("slack message structure doesn't fit")
	}

	f := &clients.Filter{}
	for _, msgAttachment := range message.OriginalMessage.Attachments {
		if f.ClusterFilterFromText(msgAttachment.Text) != nil ||
			f.AlertnameFilterFromText(msgAttachment.Text) != nil {
			return errors.New("slack message parsing for alertname and cluster failed")
		}

		incident, err := a.pdClient.GetIncident(f)
		if err != nil {
			return err
		}

		if incident.Status == clients.IncidentStatusTriggered {
			_, err = a.pdClient.AcknowledgeIncident(incident.Id, user)
			if err != nil {
				return err
			}
		}

		if user.ID == a.pdClient.GetDefaultUser().ID {
			_, err = a.pdClient.AddActualAcknowledgerAsNoteToIncident(incident.Id, slackUser.Name)
			return err
		}
	}
	return nil
}
