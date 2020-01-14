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
	"errors"
	"regexp"
)

const clusterRegex = `[\w-]*\w{2}-\w{2}-\d|admin|staging`

// ParseClusterFromString returns the cluster names found in the given string or an error.
func ParseClusterFromString(theString string) ([]string, error) {
	r, err := regexp.Compile(clusterRegex)
	if err != nil {
		return nil, err
	}

	clusters := r.FindAllString(theString, -1)
	clusters = NormalizeStringSlice(clusters)
	clusters = RemoveDuplicates(clusters)

	if len(clusters) == 0 {
		return nil, errors.New("no cluster found in input")
	}

	return clusters, nil
}
