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

package models

import (
	"github.com/nlopes/slack"
)

func appendTextSectionBlock(blocks []slack.Block, textSegments ...string) []slack.Block {
	txtBlocks := make([]*slack.TextBlockObject, 0)

	for _, txt := range textSegments {
		txtBlocks = append(txtBlocks, slack.NewTextBlockObject(
			slack.MarkdownType,
			txt,
			false,
			false,
		))
	}

	return append(blocks, slack.NewSectionBlock(nil, txtBlocks, nil))
}

func appendActionSectionBlock(blocks []slack.Block, actions ...*incidentAction) []slack.Block {
	actionBlock := slack.NewActionBlock("")
	for _, act := range actions {
		btnTxt := slack.NewTextBlockObject(slack.PlainTextType, act.text, true, false)
		actionBlock.Elements.ElementSet = append(actionBlock.Elements.ElementSet, slack.NewButtonBlockElement(act.id, act.value, btnTxt))
	}
	return append(blocks, actionBlock)
}
