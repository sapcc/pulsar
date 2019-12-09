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
	"os"

	"github.com/go-kit/kit/log"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/config"
)

// SlackClient ...
type SlackClient struct {
	logger log.Logger
	cfg    *config.SlackConfig
	client *slack.Client
}

// NewSlackClient returns a new SlackClient or an error.
func NewSlackClient(cfg *config.SlackConfig, logger log.Logger) (*SlackClient, error) {
	slackClient := slack.New(cfg.BotToken)
	if slackClient == nil {
		return nil, errors.New("failed to initialize slack client")
	}

	return &SlackClient{
		cfg:    cfg,
		logger: log.With(logger, "component", "slack"),
		client: slackClient,
	}, nil
}

// NewSlackClientFromEnv get's the configuration from the environment and returns a new SlackClient or an error.
func NewSlackClientFromEnv() (*SlackClient, error) {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg, err := config.NewSlackConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewSlackClient(cfg, logger)
}

// NewRTM returns a new RTM client.
func (s *SlackClient) NewRTM(options ...slack.RTMOption) *slack.RTM {
	return s.client.NewRTM(options...)
}

// PostMessage posts a message to the specified channel.
func (s *SlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return s.client.PostMessage(channelID, options...)
}

// GetUserByEmail returns the user or an error.
func (s *SlackClient) GetUserByEmail(email string) (*slack.User, error) {
	return s.client.GetUserByEmail(email)
}
