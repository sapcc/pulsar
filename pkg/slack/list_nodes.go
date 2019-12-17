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
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/clients"
	"github.com/sapcc/pulsar/pkg/util"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &listNodesCommand{}
	})
}

type listNodesCommand struct {
	k8sClient *clients.K8sClient
}

func (l *listNodesCommand) Init() error {
	k8sClient, err := clients.NewK8sClientFromEnv()
	if err != nil {
		return err
	}
	l.k8sClient = k8sClient
	return nil
}

func (l *listNodesCommand) IsDisabled() bool {
	return false
}

func (l *listNodesCommand) Describe() string {
	return "List nodes in cluster $clusterName."
}

func (l *listNodesCommand) Keywords() []string {
	return []string{"list nodes", "show nodes"}
}

func (l *listNodesCommand) RequiredUserRole() auth.UserRole {
	return auth.UserRoles.KubernetesUser
}

func (l *listNodesCommand) Run(msg *slack.Msg) (*slack.Msg, error) {
	clusters, err := util.ParseClusterFromString(msg.Text)
	if err != nil {
		return nil, err
	}

	// Just the first cluster.
	clusterName := clusters[0]

	if err := l.k8sClient.SetContext(clusterName); err != nil {
		return nil, err
	}

	res, err := l.k8sClient.ListNodes()
	if err != nil {
		return nil, err
	}

	return &slack.Msg{
		Type: slack.MarkdownType,
		Text: fmt.Sprintf("I found the following nodes in %s:\n```\n%s\n```", clusterName, res),
	}, nil
}
