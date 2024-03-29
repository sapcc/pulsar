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

package clients

import (
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
)

const errAlreadyReacted = "already_reacted"

// SlackClient ...
type SlackClient struct {
	logger log.Logger
	cfg    *config.SlackConfig
	client *slack.Client

}

// NewSlackClient returns a new SlackClient with Bot Token or an error.
func NewSlackBotClient(cfg *config.SlackConfig, logger log.Logger) (*SlackClient, error) {
    cfg, err := config.NewSlackConfigFromEnv()
	if err != nil {
		return nil, err
	}

    slackClient := slack.New(cfg.BotToken)
	if slackClient == nil {
		return nil, errors.New("failed to initialize slack client with bot token")
	}

	return &SlackClient{
		cfg:    cfg,
		logger: log.With(logger, "component", "slack"),
		client: slackClient,
	}, nil
}
// NewSlackClient returns a new SlackClient with Access token or an error.
func NewSlackClient(cfg *config.SlackConfig, logger log.Logger) (*SlackClient, error) {
    cfg, err := config.NewSlackConfigFromEnv()
	if err != nil {
		return nil, err
	}
    slackClient := slack.New(cfg.AccessToken)
	if slackClient == nil {
		return nil, errors.New("failed to initialize slack client with access token")
	}

	return &SlackClient{
		cfg:    cfg,
		logger: log.With(logger, "component", "slack"),
		client: slackClient,
	}, nil
}

// NewSlackClientFromEnv get's the configuration from the environment and returns a new SlackClient or an error.
func NewSlackBotClientFromEnv() (*SlackClient, error) {
	cfg, err := config.NewSlackConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewSlackBotClient(cfg, util.NewLogger())
}

// NewRTM returns a new RTM client.
func (s *SlackClient) NewRTM(options ...slack.RTMOption) *slack.RTM {
	return s.client.NewRTM(options...)
}

// PostMessage posts a message to the specified channel.
func (s *SlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	opts := []slack.MsgOption{
		slack.MsgOptionUsername(s.cfg.BotID),
		slack.MsgOptionAsUser(true),
	}

	return s.client.PostMessage(channelID, append(opts, options...)...)
}

// GetUserByEmail returns the user or an error.
func (s *SlackClient) GetUserByEmail(email string) (*slack.User, error) {
	return s.client.GetUserByEmail(email)
}

func (s *SlackClient) GetUserByID(userID string) (*slack.User, error) {
	return s.client.GetUserInfo(userID)
}

// AddReactionToMessage adds a reaction emoji to an existing message.
func (s *SlackClient) AddReactionToMessage(channel, timestamp, reaction string) error {
	msgRef := slack.NewRefToMessage(channel, timestamp)
	err := s.client.AddReaction(reaction, msgRef)
	// Ignore if reaction emoji already present.
	if err != nil && !isErrAlreadyReacted(err) {
		return err
	}
	return nil
}

// gives the message history of a channel
func (s *SlackClient) GetConversationHistory(channel string) (*slack.GetConversationHistoryResponse, error) {
    history, err := s.client.GetConversationHistory(&slack.GetConversationHistoryParameters{ ChannelID: channel, Limit: s.cfg.ChannelMessageHistoryScanCount, Inclusive: true })
	if err != nil {
		return nil, err
	}
	return history, nil
}

func isErrAlreadyReacted(err error) bool {
	if err == nil {
		return false
	}
	return strings.ToLower(err.Error()) == errAlreadyReacted
}
