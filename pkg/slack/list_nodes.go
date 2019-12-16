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
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/clients"
	"github.com/sapcc/pulsar/pkg/util"
	v1 "k8s.io/api/core/v1"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &listNodesCommand{}
	})
}

type listNodesCommand struct {
	k8sClient *clients.K8sClient
}

// Init can be used for additional initialization for the command.
func (h *listNodesCommand) Init() error {
	k8sClient, err := clients.NewK8sClientFromEnv()
	if err != nil {
		return err
	}
	h.k8sClient = k8sClient
	return nil
}

// IsDisabled can be used to (temporarily) disable the command.
func (h *listNodesCommand) IsDisabled() bool {
	return false
}

// Describe returns a brief help text for the command.
func (h *listNodesCommand) Describe() string {
	return "List nodes in a cluster."
}

// Keywords returns a list of keywords triggering the command.
func (h *listNodesCommand) Keywords() []string {
	return []string{"list nodes", "show nodes"}
}

// Run takes the slack message triggering the command and returns a slack message containing the response.
func (h *listNodesCommand) Run(msg slack.Msg) (slack.Msg, error) {
	clusters, err := util.ParseClusterFromString(msg.Text)
	if err != nil {
		return slack.Msg{}, nil
	}

	// Just the first cluster.
	clusterName := clusters[0]

	if err := h.k8sClient.SetContext(clusterName); err != nil {
		return slack.Msg{}, err
	}

	nodeList, err := h.k8sClient.ListNodes()
	if err != nil {
		return slack.Msg{}, err
	}

	if len(nodeList) == 0 {
		return slack.Msg{Text: fmt.Sprintf("No nodes in cluster %s found.", clusterName)}, nil
	}

	values := [][]string{
		{"Node", "Ready"}, // Used as headers.
	}
	for _, node := range nodeList {
		values = append(values, []string{node.Name, nodeReadyString(node)})
	}

	return util.ToSlackTable(values), nil
}

func nodeReadyString(node v1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return fmt.Sprintf("%s: %s", condition.Type, condition.Status)
		}
	}
	return ""
}
