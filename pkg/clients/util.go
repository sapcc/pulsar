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
	"fmt"
	"github.com/sapcc/go-pagerduty"
	"regexp"
	"strings"
)

// regionAlertnameRegex is used to find the region and alertname from an incident text
const regionAlertnameRegex = `.*\s\[(?P<region>[\w-]*\w{2}-\w{2}-\d|admin|staging)\]\s(?P<alertname>.+?)\s\-.*`

// parseRegionAndAlertnameFromText does what it says.
// It's meant as a workaround until Fingerprints for Prometheus alerts are supported.
// Returns an error if neither alertname nor region can be found.
func parseRegionAndAlertnameFromText(summary string) (string, string, error) {
	regionAlertnameRegex := regexp.MustCompile(regionAlertnameRegex)
	matchMap := make(map[string]string)

	match := regionAlertnameRegex.FindStringSubmatch(summary)
	for i, name := range regionAlertnameRegex.SubexpNames() {
		if i > 0 && i <= len(match) {
			m := match[i]
			if name == "" {
				continue
			} else if name == "region" {
				m = strings.ToLower(m)
			}
			matchMap[name] = m
		}
	}

	region, regionOK := matchMap["region"]
	alertname, alertnameOK := matchMap["alertname"]

	if !regionOK || !alertnameOK {
		return "", "", fmt.Errorf("pagerduty incident summary doesn not contain alertname and/or region: '%s'", summary)
	}

	return normalizeString(region), normalizeString(alertname), nil
}

func normalizeString(theString string) string {
	theString = strings.ToLower(theString)
	return strings.TrimSpace(theString)
}

func containsUser(userList []*pagerduty.User, user pagerduty.APIObject) bool {
	for _, u := range userList {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}
