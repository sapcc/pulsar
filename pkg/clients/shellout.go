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
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var errNotFound = errors.New("command not found")

// Command ...
type Command struct {
	cmd         string
	defaultArgs []string
}

// NewCommand returns a new Command or an error.
func NewCommand(cmd string, defaultArgs ...string) (*Command, error) {
	c := &Command{
		cmd:         cmd,
		defaultArgs: defaultArgs,
	}

	return c, c.verify()
}

func (c *Command) verify() error {
	err := exec.Command(c.cmd, "-v").Run()
	if err !=nil && ( strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not found")) {
		return errNotFound
	}
	return nil
}

// Run starts the command, waits until execution finished and returns stdOut or the error.
func (c *Command) Run(args ...string) (string, error) {
	cmd := exec.Command(c.cmd, append(c.defaultArgs, args...)...)

	var stdErr, stdOut bytes.Buffer
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut
	if err := cmd.Start(); err != nil {
		return "", errors.Wrap(err, stdErr.String())
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.Wrap(err, stdErr.String())
	}

	return stdOut.String(), nil
}
