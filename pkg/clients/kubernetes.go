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
	"github.com/go-kit/kit/log"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
)

// K8sClient ...
type K8sClient struct {
	cfg    *config.K8sConfig
	cmd    *Command
	logger log.Logger
}

// NewK8sClient returns a new K8sClient or an error.
func NewK8sClient(cfg *config.K8sConfig, logger log.Logger) (*K8sClient, error) {
	cmd, err := NewCommand("kubectl")
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		cfg:    cfg,
		cmd:    cmd,
		logger: log.With(logger, "component", "k8sClient"),
	}, nil
}

// NewK8sClientFromEnv get's the configuration from the environment and returns a new K8sClient or an error.
func NewK8sClientFromEnv() (*K8sClient, error) {
	cfg, err := config.NewK8sConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewK8sClient(cfg, util.NewLogger())
}

func (k *K8sClient) SetContext(context string) error {
	_, err := k.cmd.Run("config", "use-context", context)
	return err
}

// ListNodes is self-explanatory.
func (k *K8sClient) ListNodes() (string, error) {
	return k.cmd.Run("get", "nodes", "-o", "wide")
}
