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

import (
	"fmt"
	"os"
	"strings"
)

const (
	authToken    = "PAGERDUTY_AUTH_TOKEN"
	defaultEmail = "PAGERDUTY_DEFAULT_EMAIL"
    filter_services = "PAGERDUTY_SERVICES_ID_LIST"
)

// PagerdutyConfig ...
type PagerdutyConfig struct {
	AuthToken string
	DefaultEmail string
    FilterServices []string
}

// NewPagerdutyConfigFromEnv returns a new PagerdutyConfig or an error.
func NewPagerdutyConfigFromEnv() (*PagerdutyConfig, error) {
	c := &PagerdutyConfig{
		AuthToken:    os.Getenv(authToken),
		DefaultEmail: os.Getenv(defaultEmail),
        FilterServices: strings.Split(os.Getenv(filter_services), ","),
	}

	return c, c.validate()
}

func (c *PagerdutyConfig) validate() error {
	if c.AuthToken == "" {
		return fmt.Errorf("missing %s", authToken)
	}

	if c.DefaultEmail == "" {
		return fmt.Errorf("missing %s", defaultEmail)
	}

	return nil
}
