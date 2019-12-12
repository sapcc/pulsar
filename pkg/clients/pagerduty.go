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
	"github.com/sapcc/pulsar/pkg/util"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/sapcc/go-pagerduty"
	"github.com/sapcc/pulsar/pkg/config"
)

const (
	statusAcknowledged = "acknowledged"
	statusTriggered    = "triggered"
	typeUserReference  = "user_reference"
)

// PagerdutyClient wraps the pagerduty client.
type PagerdutyClient struct {
	logger          log.Logger
	cfg             *config.PagerdutyConfig
	pagerdutyClient *pagerduty.Client
	defaultUser     *pagerduty.User
}

// NewPagerdutyClient returns a new PagerdutyClient or an error.
func NewPagerdutyClient(cfg *config.PagerdutyConfig, logger log.Logger) (*PagerdutyClient, error) {
	pagerdutyClient := pagerduty.NewClient(cfg.AuthToken)
	if pagerdutyClient == nil {
		return nil, errors.New("failed to initialize pagerduty client")
	}

	c := &PagerdutyClient{
		cfg:             cfg,
		logger:          log.With(logger, "component", "pagerduty"),
		pagerdutyClient: pagerdutyClient,
	}

	defaultUser, err := c.GetUserByEmail(cfg.DefaultEmail)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting default pagerduty user with email %s", cfg.DefaultEmail)
	}
	c.defaultUser = defaultUser

	return c, nil
}

func NewPagerdutyClientFromEnv() (*PagerdutyClient, error) {
	cfg, err := config.NewPagerdutyConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewPagerdutyClient(cfg, util.NewLogger())
}

// GetUserByEmail returns the pagerduty user for the given email or an error.
func (c *PagerdutyClient) GetUserByEmail(email string) (*pagerduty.User, error) {
	userList, err := c.pagerdutyClient.ListUsers(pagerduty.ListUsersOptions{Query: email})
	if err != nil {
		return nil, err
	}

	for _, user := range userList.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with email '%s' not found", email)
}

// ListIncidents returns a list of incidents matching the given filter or an error.
func (c *PagerdutyClient) ListIncidents(f *IncidentFilter) ([]pagerduty.Incident, error) {
	o := pagerduty.ListIncidentsOptions{
		Statuses: []string{statusTriggered},
		APIListObject: pagerduty.APIListObject{
			Limit: 100,
		},
	}

	incidentList, err := c.pagerdutyClient.ListIncidents(o)
	if err != nil {
		return nil, err
	}

	// Break here if we don't need to filter
	if f == nil {
		return incidentList.Incidents, nil
	}

	return f.FilterIncidents(incidentList.Incidents), nil
}

// AcknowledgeIncident sets a incident to status acknowledged and assigns the given user to it.
func (c *PagerdutyClient) AcknowledgeIncident(incidentID string, user *pagerduty.User) error {
	incident, err := c.pagerdutyClient.GetIncident(incidentID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	userAPIObject := pagerduty.APIObject{
		Type:    typeUserReference,
		ID:      user.ID,
		Summary: user.Summary,
		HTMLURL: user.HTMLURL,
		Self:    user.Self,
	}

	incident.Status = statusAcknowledged

	incident.Acknowledgements = append(incident.Acknowledgements, pagerduty.Acknowledgement{
		At:           now.String(),
		Acknowledger: userAPIObject,
	})

	incident.Assignments = append(incident.Assignments, pagerduty.Assignment{
		At:       now.String(),
		Assignee: userAPIObject,
	})

	level.Debug(c.logger).Log("msg", "acknowledging incident", "incidentID", incident.ID, "userEmail", user.Email)
	return c.pagerdutyClient.ManageIncidents(user.Email, []pagerduty.Incident{*incident})
}

// AddActualAcknowledgerAsNoteToIncident adds a note containing the actual acknowledger to the given incident.
func (c *PagerdutyClient) AddActualAcknowledgerAsNoteToIncident(incidentID, actualAcknowledger string) error {
	now := time.Now().UTC()
	note := pagerduty.IncidentNote{
		CreatedAt: now.String(),
		User: pagerduty.APIObject{
			ID:      c.defaultUser.ID,
			Type:    typeUserReference,
			Summary: c.defaultUser.Summary,
			Self:    c.defaultUser.Self,
			HTMLURL: c.defaultUser.HTMLURL,
		},
		Content: fmt.Sprintf("Incident was acknowledged on behalf of %s. time: %s", actualAcknowledger, now.String()),
	}

	return c.pagerdutyClient.CreateIncidentNote(incidentID, c.defaultUser.Email, note)
}

// ListTodaysOnCalls returns the OnCall users for today.
func (c *PagerdutyClient) ListTodaysOnCallUsers(scheduleID *string) ([]*pagerduty.User, error) {
	listOpts := pagerduty.ListOnCallOptions{}
	listOpts.Limit = 100
	listOpts.Earliest = true
	listOpts.Since = util.TimestampToString(util.TimeStartOfDay())
	listOpts.Until = util.TimestampToString(util.TimeEndOfDay())

	if scheduleID != nil {
		listOpts.Includes = []string{"schedules"}
		listOpts.ScheduleIDs = []string{*scheduleID}
	}

	onCallList, err := c.pagerdutyClient.ListOnCalls(listOpts)
	if err != nil {
		return nil, err
	}

	// Deduplicate.
	res := make([]*pagerduty.User, 0)
	for _, onCall := range onCallList.OnCalls {
		if !containsUser(res, onCall.User) {
			u, err := c.pagerdutyClient.GetUser(onCall.User.ID, pagerduty.GetUserOptions{})
			if err != nil {
				level.Error(c.logger).Log("msg", "error getting user", "id", onCall.User.ID, "err", err.Error())
				continue
			}

			res = append(res, u)
		}
	}

	return res, nil
}

// GetSchedule returns a pagerduty schedule for the given name or an error.
func (c *PagerdutyClient) GetSchedule(scheduleName string) (*pagerduty.Schedule, error) {
	listOpts := pagerduty.ListSchedulesOptions{}
	listOpts.Limit = 100
	listOpts.Query = scheduleName

	scheduleList, err := c.pagerdutyClient.ListSchedules(listOpts)
	if err != nil {
		return nil, err
	}

	for _, sched := range scheduleList.Schedules {
		if normalizeString(sched.Name) == normalizeString(scheduleName) {
			return &sched, nil
		}
	}

	return nil, fmt.Errorf("schedule not found")
}
