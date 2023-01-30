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

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	botToken                  = "SLACK_BOT_TOKEN"
	botID                     = "SLACK_BOT_ID"
	authorizedUserGroupNames  = "SLACK_AUTHORIZED_USER_GROUP_NAMES"
	kubernetesUserGroupNames  = "SLACK_KUBERNETES_USER_GROUP_NAMES"
	kubernetesAdminGroupNames = "SLACK_KUBERNETES_ADMIN_GROUP_NAMES"
	accessToken               = "SLACK_ACCESS_TOKEN"
	verificationToken         = "SLACK_VERIFICATION_TOKEN"
    channelIdsListForPdSync   = "SLACK_CHANNELS_ID_LIST"
    channelMessageHistoryScanCount = "SLACK_CHANNELS_MESSAGE_HISTORY_SCAN_COUNT"
	apiPort                   = "API_PORT"
	apiHost                   = "API_HOST"
)

// SlackConfig ...
type SlackConfig struct {
	// BotToken is the Slack token with bot permissions.
	BotToken string

	// BotID is the id of the bot.
	BotID string

	// AccessToken is the Slack token with permissions to list user groups and their members.
	AccessToken string

	// VerificationToken used to verify messages from the Slack API.
	VerificationToken string

	// AuthorizedUserGroupNames is the list of user group names whose members are authorized to interact with the bot.
	AuthorizedUserGroupNames []string

	// KubernetesUserGroupNames is the list of user group names whose members are authorized to perform read operations for kubernetes clusters via the bot.
	KubernetesUserGroupNames []string

	// KubernetesAdminGroupNames is the list of user group names whose members are authorized to perform all operations for kubernetes clusters via the bot.
	KubernetesAdminGroupNames []string

	// APIPort is the port on which the API is exposed.
	APIPort int

	// APIHost is the host on which the API is exposed.
	APIHost string

    // Slack Channel Ids for PD incident sync
    ChannelIdsListForPdSync []string

    // Slack Channel History Message Count which will be scanned 
    ChannelMessageHistoryScanCount int
}

func NewSlackConfigFromEnv() (*SlackConfig, error) {
	port := 8080
	if p, err := strconv.Atoi(os.Getenv(apiPort)); err == nil {
		port = p
	}

	host := "0.0.0.0"
	if h := os.Getenv(apiHost); h != "" {
		host = h
	}

    defaultChannelMessageHistoryScanCount := 20
    if msc, err := strconv.Atoi(os.Getenv(channelMessageHistoryScanCount)); err == nil {
        defaultChannelMessageHistoryScanCount = msc
    }

	c := &SlackConfig{
		BotToken:                  os.Getenv(botToken),
		BotID:                     os.Getenv(botID),
		AccessToken:               os.Getenv(accessToken),
		VerificationToken:         os.Getenv(verificationToken),
        ChannelIdsListForPdSync:   strings.Split(os.Getenv(channelIdsListForPdSync), ","),
        ChannelMessageHistoryScanCount: defaultChannelMessageHistoryScanCount,
		AuthorizedUserGroupNames:  strings.Split(os.Getenv(authorizedUserGroupNames), ","),
		KubernetesUserGroupNames:  strings.Split(os.Getenv(kubernetesUserGroupNames), ","),
		KubernetesAdminGroupNames: strings.Split(os.Getenv(kubernetesAdminGroupNames), ","),
		APIHost:                   host,
		APIPort:                   port,
	}
	return c, c.validate()
}

func (c *SlackConfig) validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("missing %s", botToken)
	}
	if c.BotID == "" {
		return fmt.Errorf("missing %s", botID)
	}
	if len(c.AuthorizedUserGroupNames) == 0 {
		return fmt.Errorf("missing %s", authorizedUserGroupNames)
	}
	if c.AccessToken == "" {
		return fmt.Errorf("missing %s", accessToken)
	}
	if c.VerificationToken == "" {
		return fmt.Errorf("missing %s", verificationToken)
	}
    if len(c.ChannelIdsListForPdSync) == 0 {
        return fmt.Errorf("missing or empty %s", verificationToken)
    }

	return nil
}
