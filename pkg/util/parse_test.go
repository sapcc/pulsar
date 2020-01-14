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

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseClusterFromString(t *testing.T) {
	stimuli := map[string][]string{
		"[CRITICAL - 2] [LA-BR-1] ManyPodsNotReadyOnNode - Less then 75% of pods ready on node": {"la-br-1"},
		"[RESOLVED] [LA-BR-1] OpenstackNovaDatapathDown - Datapath nova metadata is down": {"la-br-1"},
		"[CRITICAL] [S-LA-BR-1] InfrastructurePrometheusFederationFailed - Infrastructure Prometheus s-la-br-1 is down": {"s-la-br-1"},
	}

	for inputString, expected := range stimuli {
		got, err := ParseClusterFromString(inputString)
		assert.NoError(t, err, "there should be no error parsing cluster names from the string")
		assert.EqualValues(t, expected, got, "result and expected should have equal values")
	}
}
