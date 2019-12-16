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
	"math/rand"

	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/bot"
)

// init registers the command.
func init() {
	bot.RegisterCommand(func() bot.Command {
		return &helloCommand{}
	})
}

type helloCommand struct{}

// Init can be used for additional initialization for the command.
func (h *helloCommand) Init() error {
	return nil
}

// IsDisabled can be used to (temporarily) disable the command.
func (h *helloCommand) IsDisabled() bool {
	return false
}

// Describe returns a brief help text for the command.
func (h *helloCommand) Describe() string {
	return "Be polite. Say hello."
}

// Keywords returns a list of keywords triggering the command.
func (h *helloCommand) Keywords() []string {
	return []string{"hey", "hello", "hi"}
}

// Run takes the slack message triggering the command and returns a slack message containing the response.
func (h *helloCommand) Run(msg *slack.Msg) (*slack.Msg, error) {

	greetings := []string{
		fmt.Sprintf("What's up <@%s>?", msg.User),
		"How may I be of service?",
		"Hello :ccloud:!",
		fmt.Sprintf("Hello <@%s>.", msg.User),
		fmt.Sprintf("Nice to meet you <@%s>.", msg.User),
	}

	return &slack.Msg{Text: greetings[rand.Intn(len(greetings))]}, nil
}
