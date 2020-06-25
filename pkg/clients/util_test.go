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
	"testing"

	"github.com/sapcc/pulsar/pkg/util"
	"github.com/stretchr/testify/assert"
)

const (
	slackText          				= "\n*[CRITICAL]* *[ap-sa-1]* VCenterRedundancyLostHAPolicyFaulty - VC vc-b-0.cc.eu-nl-1.cloud.sap has a faulty AdmissionControlPolicy for cluster XYZ, failover will not work.\n:fire: VC ... \n"
	slackTextWithLink 				= "\n*[CRITICAL]* *[AP-JP-1]* *<https://alertmanager.somewhere.cloud.com/#/alerts?receiver=slack_api_critical|OpenstackDatapathDown>* - Blackbox datapath test\n:fire: Datapath maia_metrics is down for 15 times in a row. ... \n"
	slackTextMulti					= "\n*[CRITICAL - 6]* *[EU-RU-1]* NetworkApicProcessMaxMemoryUsedCritical - \n:fire: Max memory 2.817073152e+09 used by process nfm/topology/pod-1/node-000/sys/procsys/proc-10560 on apic host .. \n"
	slackTextMultipleNoDescription 	= "\n*[CRITICAL - 6]* *[EU-RU-1]* NetworkApicProcessMaxMemoryUsedCritical - "
	)

func TestParseAlertFromSlackMessageText(t *testing.T) {
	// mapping of input string to expected result map
	tests := map[string]map[string]string{
		slackText: {
			"alertname": "VCenterRedundancyLostHAPolicyFaulty",
			"region":    "ap-sa-1",
		},
		slackTextWithLink: {
			"alertname": "OpenstackDatapathDown",
			"region":    "AP-JP-1",
		},
		slackTextMulti: {
			"alertname": "NetworkApicProcessMaxMemoryUsedCritical",
			"region":    "EU-RU-1",
		},
		slackTextMultipleNoDescription: {
			"alertname": "NetworkApicProcessMaxMemoryUsedCritical",
			"region":    "EU-RU-1",
		},
	}

	for stimuli, expectedMap := range tests {
		region, alertname, err := parseRegionAndAlertnameFromText(stimuli)
		assert.NoError(t, err, "there should be no error parsing the slack message text: %s", stimuli)
		assert.Equal(t, util.NormalizeString(expectedMap["alertname"]), alertname, "the alertname should be equal")
		assert.Equal(t, util.NormalizeString(expectedMap["region"]), region, "the region should be equal")
	}
}
