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

import (
	"fmt"
	"time"

	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/util"
)

const (
	emojiGreenCheckmark = ":green_checkmark:"
	emojiRotatingLight  = ":rotating_light:"
	emojiTV             = ":tv:"
	emojiNotebook       = ":notebook:"
	emojiClock          = ":clock1:"

	statusOpen   = "open"
	statusClosed = "closed"
)

type Incident struct {
	title,
	description,
	statusEmoji string
	startTime,
	endTime time.Time
	duration time.Duration
	isClosed bool
	reporter,
	lead *User
	severity Severity
}

func NewIncident(title string, reporter, lead *User, severity Severity) *Incident {
	return &Incident{
		title:     title,
		reporter:  reporter,
		lead:      lead,
		severity:  severity,
		startTime: time.Now().UTC(),
	}
}

func (i *Incident) SetLead(lead *User) {
	i.lead = lead
}

func (i *Incident) SetDescription(description string) {
	i.description = description
}

func (i *Incident) Close() {
	i.isClosed = true
	i.endTime = time.Now().UTC()
	i.duration = i.endTime.Sub(i.startTime)
	i.statusEmoji = emojiGreenCheckmark
}

func (i *Incident) ToSlackMessage() *slack.Msg {
	blocks := make([]slack.Block, 0)

	// Header block with title and involved users.
	blocks = appendTextSectionBlock(blocks, fmt.Sprintf("*Incident*: %s", i.title))
	blocks = appendTextSectionBlock(blocks, i.reporter.String())
	blocks = appendTextSectionBlock(blocks, i.lead.String())

	blocks = append(blocks, slack.NewDividerBlock())

	// Block with incident details.
	status := statusOpen
	if i.isClosed {
		status = statusClosed
	}
	blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Status: %s", emojiTV, status))
	blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Severity: %s", emojiRotatingLight, i.severity.String()))
	blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Started: %s", emojiClock, util.HumanizeTimestamp(i.startTime)))

	if i.isClosed {
		blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Closed: %s", emojiClock, util.HumanizeTimestamp(i.endTime)))
		blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Duration: %s", emojiClock, util.HumanizeDuration(i.duration)))
	}

	if i.description != "" {
		blocks = appendTextSectionBlock(blocks, fmt.Sprintf("%s Description: %s", emojiNotebook, i.description))
	}

	blocks = append(blocks, slack.NewDividerBlock())

	// Footer block for buttons.
	blocks = appendActionSectionBlock(blocks,
		newIncidentAction("closeID", "Close", "close"),
		newIncidentAction("editID", "Edit", "edit"),
		newIncidentAction("pageID", "Page on-call", "page"),
	)

	blockMsg := slack.NewBlockMessage(blocks...)
	return &blockMsg.Msg
}
