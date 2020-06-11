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

package cmd

import (
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/api"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/bot"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
	"github.com/sapcc/pulsar/pkg/version"
	"github.com/spf13/cobra"

	// Load all slack plugins.
	_ "github.com/sapcc/pulsar/pkg/slack"
)

const rootCmdLongUsage = "Pulsar bot mode"

func New() *cobra.Command {
	stop := make(chan struct{})

	cmd := &cobra.Command{
		Use:          "pulsar",
		Short:        "Slack bot mode",
		Long:         rootCmdLongUsage,
		SilenceUsage: true,
		Version:      version.Print(),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := util.NewLogger()

			cfg, err := config.NewSlackConfigFromEnv()
			if err != nil {
				return err
			}

			authorizer, err := auth.New(cfg, logger)
			if err != nil {
				return errors.Wrap(err, "error initializing authorizer")
			}

			// Start the bot.
			b, err := bot.New(authorizer, cfg, logger)
			if err != nil {
				return errors.Wrap(err, "error initializing bot")
			}

			// Start the API handling interactive messages.
			a, err := api.New(authorizer, cfg, logger)
			if err != nil {
				return errors.Wrap(err, "error initializing api")
			}

			go authorizer.Run(stop)
			go a.Serve(stop)
			go b.ListenAndRespond(stop)

			<-stop
			close(stop)

			return nil
		},
	}

	return cmd
}
