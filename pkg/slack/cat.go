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

package slack

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/bot"
)

func init() {
	bot.RegisterCommand(func() bot.Command {
		return &catCommand{}
	})
}

type catCommand struct{}

func (c *catCommand) Init() error {
	return nil
}

func (c *catCommand) IsDisabled() bool {
	return false
}

func (c *catCommand) Describe() string {
	return "Post cute cat pics"
}

func (c *catCommand) Keywords() []string {
	return []string{"cat", "kitty"}
}

func (c *catCommand) Run(msg *slack.Msg) (*slack.Msg, error) {
	resp, err := http.Get("https://api.thecatapi.com/api/images/get?format=xml&size=med&results_per_page=1")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received reponse with status code: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}
	defer resp.Body.Close()

	var res struct {
		Data struct {
			Images []struct {
				Image struct {
					URL string `xml:"url"`
				} `xml:"image"`
			} `xml:"images"`
		} `xml:"data"`
	}

	if err := xml.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	imgs := res.Data.Images
	if len(imgs) == 0 {
		return nil, errors.New("no cat picture found")
	}

	imgURL := imgs[0].Image.URL
	if imgURL == "" {
		return nil, errors.New("no cat picture found")
	}

	return &slack.Msg{
		Text: fmt.Sprintf("Here's a cute cat pic for you:\n%s", imgURL),
	}, nil
}
