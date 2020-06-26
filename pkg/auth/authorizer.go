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

package auth

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
)

// Authorizer ...
type Authorizer struct {
	logger         log.Logger
	cfg            *config.SlackConfig
	client         *slack.Client
	tickerInterval time.Duration

	authorizedUserIDs,
	kubernetesAdminsUserIDs,
	kubernetesUsersUserIDs []string
}

// New returns a new Authorizer or an error.
func New(cfg *config.SlackConfig, logger log.Logger) (*Authorizer, error) {
	c := slack.New(cfg.AccessToken)
	if c == nil {
		return nil, errors.New("cannot create slack client")
	}

	a := &Authorizer{
		logger:         log.With(logger, "component", "authorizer"),
		cfg:            cfg,
		client:         c,
		tickerInterval: 10 * time.Minute,
	}

	if err := a.getAuthorizedUserIDs(); err != nil {
		return nil, err
	}

	return a, nil
}

// IsUserAuthorized checks whether the given user is authorized to run the bot command.
func (a *Authorizer) IsUserAuthorized(userID string, requiredUserRole UserRole) bool {
	switch requiredUserRole {
	case UserRoles.Base:
		return util.Contains(a.authorizedUserIDs, userID)
	case UserRoles.KubernetesUser:
		return util.Contains(a.kubernetesUsersUserIDs, userID)
	case UserRoles.KubernetesAdmin:
		return util.Contains(a.kubernetesAdminsUserIDs, userID)
	}

	return false
}

// Run starts the continuous synchronization in the background.
func (a *Authorizer) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(a.tickerInterval)
	defer ticker.Stop()

	go func() {
		for{
			select {
			case <-ticker.C:
				if err := a.getAuthorizedUserIDs(); err != nil {
					level.Error(a.logger).Log("msg", "failed to refresh authorized users", "err", err.Error())
				}
			case <-stop:
				return
			}
		}
	}()

	<-stop
}

func (a *Authorizer) getAuthorizedUserIDs() error {
	ugList, err := a.client.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return errors.Wrap(err, "failed to list user groups")
	}

	for _, ug := range ugList {
		if util.Contains(a.cfg.AuthorizedUserGroupNames, ug.Name) {
			a.authorizedUserIDs = append(a.authorizedUserIDs, ug.Users...)
		}

		if util.Contains(a.cfg.KubernetesUserGroupNames, ug.Name) {
			a.kubernetesUsersUserIDs = append(a.kubernetesUsersUserIDs, ug.Users...)
		}

		if util.Contains(a.cfg.KubernetesAdminGroupNames, ug.Name) {
			a.kubernetesAdminsUserIDs = append(a.kubernetesAdminsUserIDs, ug.Users...)
		}
	}

	if a.authorizedUserIDs == nil || len(a.authorizedUserIDs) == 0 {
		return errors.New("not a single user is authorized to respond to slack messages. check configured authorized user groups")
	}

	level.Debug(a.logger).Log("msg", "syncing authorized users from slack groups")
	return nil
}
