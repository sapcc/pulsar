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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // https://github.com/kubernetes/client-go/issues/242
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient ...
type K8sClient struct {
	cfg       *config.K8sConfig
	clientset *kubernetes.Clientset
	logger    log.Logger
}

// NewK8sClient returns a new K8sClient or an error.
func NewK8sClient(cfg *config.K8sConfig, logger log.Logger) (*K8sClient, error) {
	clientset, err := newClientSetWithContext("", cfg.KubeConfig)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		cfg:       cfg,
		clientset: clientset,
		logger:    log.With(logger, "component", "k8sClient"),
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

// SetContext sets the kube context.
func (k *K8sClient) SetContext(context string) error {
	clientset, err := newClientSetWithContext(context, k.cfg.KubeConfig)
	if err != nil {
		return err
	}

	k.clientset = clientset
	return nil
}

// ListNodes is self-explanatory.
func (k *K8sClient) ListNodes() ([]v1.Node, error) {
	nodeList, err := k.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

func newClientSetWithContext(context, kubeConfigPath string) (*kubernetes.Clientset, error) {
	// Set the path to the kubeConfig if given.
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeConfigPath != "" {
		rules.ExplicitPath = kubeConfigPath
	}

	// Set the context if given.
	overrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		overrides.CurrentContext = context
	}

	kubeConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}
