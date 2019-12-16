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

import "os"

// K8sConfig ...
type K8sConfig struct {
	// Path to kubeconfig.
	KubeConfig string
}

// NewK8sConfigFromEnv returns a new K8sConfig or an error.
func NewK8sConfigFromEnv() (*K8sConfig, error) {
	k := &K8sConfig{
		KubeConfig: os.Getenv("KUBECONFIG"),
	}
	return k, k.validate()
}

func (k *K8sConfig) validate() error {
	return nil
}
