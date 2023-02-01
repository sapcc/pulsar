/*******************************************************************************
*
* Copyright 2023 SAP SE
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
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/go-kit/log/level"
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/clients"
)

// incident sync will do:
// run frequently as cron job
// try to match open incidents from pagerduty of defined services to defined slack channels
// filter on service from environmental values: PD_SERVICES_ID_LIST
// filter on slack channels from environmental values: SLACK_CHANNELS_ID_LIST
func (a *API) pd_slack_incidents_sync() error {

	// 1. get all PD incidents from our defined services
	f := &clients.Filter{}
	f.SetLimit(100)
	iL, err := a.pdClient.ListIncidents(f)
	if err != nil {
		return err
	}
	for _, i := range iL {
		a.enrich_slack_channel_with_incident(&i)
	}

	return nil
}

func (a *API) enrich_slack_channel_with_incident(incident *pagerduty.Incident) {
	level.Info(a.logger).Log("write into Channel")

	for _, channelId := range a.cfg.ChannelIdsListForPdSync {

		h, err := a.slackClient.GetConversationHistory(channelId)
		if err != nil {
			level.Error(a.logger).Log("trouble to get slack channel messages - ", channelId, ": ", err)
			continue
		}
		for _, message := range h.Messages {

            // skip if message isn't in the same time range (other alert)
            if !a.checkIfIncidentMessageTimeIsMoreOrLessSame(&message, incident){
                level.Debug(a.logger).Log(fmt.Printf("skip message in channel %s because of timestamp.", message.Channel))
                continue
            }

            level.Debug(a.logger).Log("Incident Match in Channel ", channelId, " for incident ", incident.APIObject.HTMLURL)

            // AlertManager makes Attachment Messages
			if len(message.Attachments) > 0 {

                s := strings.ToLower(message.Attachments[0].Text)
				level.Info(a.logger).Log(s)
				region, alertname, err := clients.ParseRegionAndAlertnameFromText(incident.Summary)

                // skip resolved
                if err == nil && strings.Contains(s, "resolved") {
                    level.Debug(a.logger).Log(fmt.Printf("skip message in channel %s because it's resolved", message.Channel))
                    continue
                }

                // get state (acknowledged / already add PD link)
                bIconFireFighter, bIconPagerduty := a.checkReactions(&message)
                if bIconFireFighter && bIconPagerduty {
                    level.Debug(a.logger).Log(fmt.Printf("skip message in channel %s because it's already handled", message.Channel))
                    continue
                }

                // handle - both states (acknowledged / pd link added) are possible
				if err == nil && strings.Contains(strings.ToLower(s), region) && strings.Contains(strings.ToLower(s), alertname) {

                    message.Channel = channelId

                    // add pd incident link add
                    if !bIconPagerduty {
                        level.Debug(a.logger).Log(fmt.Printf("add pd link to slack message - channel %s for alertname: %s", channelId, alertname))
                        a.addPdLink(&message, incident)
                        a.addReactionHandled(&message)
                    }

                    // if not marked as acknowledged - we mark it
					if incident.Status == "acknowledged" && !bIconFireFighter {
                            level.Debug(a.logger).Log(fmt.Printf("add ack Icon to slack message - channel %s for alertname: %s", channelId, alertname))
                            a.addReactionAcknowledged(&message, incident)
                    }

				}
			}
		}
    }
}

func (a *API)checkIfIncidentMessageTimeIsMoreOrLessSame(message *slack.Message, incident *pagerduty.Incident) bool {
    tp, _ := time.Parse(time.RFC3339, incident.CreatedAt)
    tmm, err := strconv.ParseInt(strings.Split(message.Timestamp, ".")[0], 10, 64)
    if err != nil {
        level.Error(a.logger).Log(err)
    }
    return tp.Sub(time.Unix(tmm,0)).Abs().Minutes() <= 1
}

func (a *API)checkReactions(message *slack.Message) (bool, bool) {

    bIconFireFighter :=false
    bIconPagerduty := false

    for _, r := range message.Reactions{
        if r.Name == emojiPagerDuty {
            level.Debug(a.logger).Log(fmt.Printf("there is a reaction icon :%s:  - we handled it already - channel: %s.", emojiPagerDuty, message.Channel))
            bIconPagerduty = true
        }

        if r.Name == emojiFirefighter {
            level.Debug(a.logger).Log(fmt.Printf("there is a reaction icon :%s:  - we handled it already - channel: %s.", emojiFirefighter, message.Channel))
            bIconFireFighter = true
        }
    }

    return bIconFireFighter, bIconPagerduty
}

// Post incident link
func (a *API)addPdLink(message *slack.Message, incident *pagerduty.Incident) error {
    if _, _, err := a.slackBotClient.PostMessage(
        message.Channel,// channelId,
        slack.MsgOptionText(fmt.Sprintf("PD Incident (%d): %s", incident.IncidentNumber, incident.HTMLURL), false),
        slack.MsgOptionTS(message.Timestamp),
    ); err != nil {
        level.Error(a.logger).Log(err)
        return err
    }
    return nil
}

// Add reaction emoji to original message.
func (a *API)addReactionHandled(message *slack.Message) error {
    if err := a.slackBotClient.AddReactionToMessage(
        message.Channel,
        message.Timestamp,
        emojiPagerDuty,
    ); err != nil {
        level.Error(a.logger).Log(err)
        return err
    }
    return nil
}

func (a *API)addReactionAcknowledged(message *slack.Message, incident *pagerduty.Incident) error {
    // Add reaction emoji to original message.
    if err := a.slackBotClient.AddReactionToMessage(
        message.Channel,
        message.Timestamp,
        emojiFirefighter,
    ); err != nil {
        level.Error(a.logger).Log(err)
        return err
    }
    // Post the message.
    if _, _, err := a.slackBotClient.PostMessage(
        message.Channel,
        slack.MsgOptionText(fmt.Sprintf(acknowledgeString, incident.Acknowledgements[0].Acknowledger.Summary), false),
        slack.MsgOptionTS(message.Timestamp),
    ); err != nil {
        level.Error(a.logger).Log(err)
        return err
    }
    return nil
}
