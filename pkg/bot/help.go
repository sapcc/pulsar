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

package bot

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/util"
)

// Help is the only command available from the beginning.
// Other commands shall be implemented following the factory pattern in the slack package.
type helpCommand struct {
	availableCommands []Command
}

func (b *Bot) newHelpCommand(availableCommands []Command) Command {
	return &helpCommand{
		availableCommands: availableCommands,
	}
}

func (h *helpCommand) Init() error {
	return nil
}

func (h *helpCommand) IsDisabled() bool {
	return false
}

func (h *helpCommand) Keywords() []string {
	return []string{"help"}
}

func (h *helpCommand) Describe() string {
	return "Help for all commands"
}

func (h *helpCommand) Run(msg slack.Msg) (slack.Msg, error) {
	commands := [][]string{
		{"*Command*", "*Description*"},
	}
	for _, c := range h.availableCommands {
		commands = append(commands, []string{strings.Join(c.Keywords(), ", "), c.Describe()})
	}

	return util.ToSlackTable(commands), nil
}
