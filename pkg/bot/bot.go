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
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/clients"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
)

// Bot is the struct for the slack bot.
type Bot struct {
	authorizer  *auth.Authorizer
	logger      log.Logger
	client      *clients.SlackClient
	rtmClient   *slack.RTM
	botID       string
	channelID   string
	helpCommand Command
	commands    []Command
}

// New returns a new Bot or an error.
func New(authorizer *auth.Authorizer, cfg *config.SlackConfig, logger log.Logger) (*Bot, error) {
	slackClient, err := clients.NewSlackClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		authorizer: authorizer,
		logger:     log.With(logger, "component", "bot"),
		client:     slackClient,
		rtmClient:  slackClient.NewRTM(),
		botID:      cfg.BotID,
	}

	for _, c := range availableCommands {
		cmd := c()
		if err := cmd.Init(); err != nil {
			level.Error(b.logger).Log("msg", "failed to initialize command", "keywords", strings.Join(cmd.Keywords(), ", "), "description", cmd.Describe(), "err", err.Error())
			continue
		}
		level.Debug(b.logger).Log("msg", "registering command", "keywords", strings.Join(cmd.Keywords(), ", "), "description", cmd.Describe())
		b.commands = append(b.commands, cmd)
	}

	b.helpCommand = b.newHelpCommand(b.commands)
	b.commands = append(b.commands, b.helpCommand)

	return b, nil
}

// ListenAndRespond will make the bot listen to events and respond o them.
func (b *Bot) ListenAndRespond(stop <-chan struct{}) {
	// Listen to slack events.
	go b.rtmClient.ManageConnection()

	for {
		select {
		case msg := <-b.rtmClient.IncomingEvents:
			level.Debug(b.logger).Log("msg", "received slack event", "type", msg.Type)

			switch e := msg.Data.(type) {
			case *slack.MessageEvent:
				if err := b.handleMessageEvent(e); err != nil {
					level.Error(b.logger).Log("msg", "error handling slack event", "err", err.Error())
					b.respond(&slack.Msg{Text: "Failed to respond"}, &e.Msg)
				}

			case *slack.RTMError:
				level.Error(b.logger).Log("msg", "slack RTM error", "err", e.Error())

			case *slack.InvalidAuthEvent:
				level.Error(b.logger).Log("msg", "slack authentication failed")

			case *slack.ConnectionErrorEvent:
				level.Error(b.logger).Log("error connecting to slack", "err", e.Error())
			}
		}
	}
}

func (b *Bot) handleMessageEvent(e *slack.MessageEvent) error {
	info := b.rtmClient.GetInfo()
	prefix := fmt.Sprintf("<@%s>", info.User.ID)

	if !strings.HasPrefix(e.Text, prefix) {
		return nil
	}

	// Only respond if the bot is mentioned.
	text := e.Msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	// Update original message text with normalized one.
	e.Msg.Text = text

	atLeastOneCommand := false
	for _, c := range b.commands {
		if util.HasAnyPrefix(c.Keywords(), text) {

			if !b.authorizer.IsUserAuthorized(e.Msg.User, c.RequiredUserRole()) {
				level.Debug(b.logger).Log("msg", "user is not authorized", "userID", e.Msg.User, "requiredRole", c.RequiredUserRole())
				b.respond(&slack.Msg{Text: "You are not authorized :x:"}, &e.Msg)
				return nil
			}

			atLeastOneCommand = true
			response, err := c.Run(&e.Msg)
			if err != nil {
				return err
			}

			if err := b.respond(response, &e.Msg); err != nil {
				return err
			}
		}
	}

	if atLeastOneCommand {
		return nil
	}

	response, err := b.helpCommand.Run(&e.Msg)
	if err != nil {
		return err
	}

	return b.respond(response, &e.Msg)
}

func (b *Bot) respond(msg, originalMsg *slack.Msg) error {
	opts := []slack.MsgOption{
		slack.MsgOptionUsername(b.botID),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText(msg.Text, false),
	}

	if len(msg.Blocks.BlockSet) > 0 {
		opts = append(opts, slack.MsgOptionBlocks(msg.Blocks.BlockSet...))
	}

	if msg.Attachments != nil {
		opts = append(opts, slack.MsgOptionAttachments(msg.Attachments...))
	}

	_, _, err := b.client.PostMessage(originalMsg.Channel, opts...)
	return err
}
