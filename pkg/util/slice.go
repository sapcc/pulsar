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
	"strings"
)

// Contains checks whether the given string slice contains the searchString.
func Contains(sslice []string, searchString string) bool {
	for _, s := range sslice {
		if s == searchString {
			return true
		}
	}
	return false
}

// HasAnyPrefix checks whether the given string starts with any of the prefixes.
func HasAnyPrefix(prefixes []string, theString string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(theString, p) {
			return true
		}
	}
	return false
}

// TrimAnyPrefix trims any of the given prefixes from the string
func TrimAnyPrefix(prefixes []string, theString string) string {
	for _, p := range prefixes {
		theString = strings.TrimPrefix(theString, p)
	}
	return theString
}

// IsSlicesEqual check whether the given slice have equal content but not necessarily in the same order.
func IsSlicesEqual(sslice1, sslice2 []string) bool {
	if len(sslice1) != len(sslice2) {
		return false
	}

	visited := make([]bool, len(sslice2))
	for i := 0; i < len(sslice1); i++ {
		found := false
		for j := 0; j < len(sslice2); j++ {
			if visited[j] {
				continue
			}
			if sslice1[i] == sslice2[j] {
				visited[j] = true
				found = true
			}
		}

		if !found {
			return false
		}

	}

	return true
}

// RemoveDuplicates does what it says on the given string slice.
func RemoveDuplicates(stringSlice []string) []string {
	keys := make(map[string]bool)
	res := make([]string, 0)
	for _, itm := range stringSlice {
		if _, exists := keys[itm]; !exists {
			keys[itm] = true
			res = append(res, itm)
		}
	}
	return res
}
