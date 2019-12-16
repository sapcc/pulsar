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

import "github.com/nlopes/slack"

// ToSlackTable create a 2 column slack table using the TextBlockObjects.
// The input shall look like:
// [][]string{
//   {"headerColumn1", "headerColumn2"},
//	 {"valueLine1", "valueLine2"},
// }
func ToSlackTable(values [][]string) *slack.Msg {
	fields := make([]*slack.TextBlockObject, 0)
	for _, v := range values {
		for _, itm := range v {
			fields = append(fields, slack.NewTextBlockObject(slack.PlainTextType, itm, true, false))
		}
	}

	blocks := slack.Blocks{}
	blocks.BlockSet = append(blocks.BlockSet, slack.NewSectionBlock(nil, fields, nil))
	return &slack.Msg{Blocks: blocks}
}
