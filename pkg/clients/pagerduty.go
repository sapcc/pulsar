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
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/sapcc/pulsar/pkg/config"
	"github.com/sapcc/pulsar/pkg/util"
)

const (
	IncidentStatusAcknowledged = "acknowledged"
	IncidentStatusTriggered    = "triggered"
	typeUserReference          = "user_reference"
	typeIncident               = "incident"
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

// GetDefaultUser returns the pagerduty default user.
func (c *PagerdutyClient) GetDefaultUser() *pagerduty.User {
	return c.defaultUser
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
func (c *PagerdutyClient) ListIncidents(f *Filter) ([]pagerduty.Incident, error) {
	o := pagerduty.ListIncidentsOptions{
		Statuses: []string{IncidentStatusTriggered, IncidentStatusAcknowledged},
		Since: time.Now().AddDate(0,0,-1).Format(time.RFC3339),
		SortBy: "created_at:desc",
	}

	if f != nil {
		level.Debug(c.logger).Log("msg", "listing pagerduty incident", "filter", f.ToString())
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

// GetIncident returns the latest incident matching the filter or an error.
func (c *PagerdutyClient) GetIncident(f *Filter) (*pagerduty.Incident, error) {
	// Return the most recent incident.
	f.SetLimit(1)
	incidentList, err := c.ListIncidents(f)
	if err != nil {
		return nil, errors.Wrap(err, "error listing pagerduty incidents")
	}

	if len(incidentList) == 0 {
		return nil, errors.New("not a single pagerduty incident found")
	}

	return &incidentList[0], nil
}

// AcknowledgeIncident sets a incident to status acknowledged and assigns the given user to it.
func (c *PagerdutyClient) AcknowledgeIncident(incidentID string, user *pagerduty.User) (*pagerduty.ListIncidentsResponse, error) {
	if user == nil {
		user = c.defaultUser
	}

	incident := pagerduty.ManageIncidentsOptions{
		ID:     incidentID,
		Type:   typeIncident,
		Status: IncidentStatusAcknowledged,
	}

	level.Debug(c.logger).Log("msg", "acknowledging incident", "incidentID", incident.ID, "userEmail", user.Email)
	return c.pagerdutyClient.ManageIncidents(user.Email, []pagerduty.ManageIncidentsOptions{incident})
}

// AddActualAcknowledgerAsNoteToIncident adds a note containing the actual acknowledger to the given incident.
func (c *PagerdutyClient) AddActualAcknowledgerAsNoteToIncident(incidentID, actualAcknowledger string) (*pagerduty.IncidentNote, error) {
	now := time.Now().UTC()
	note := pagerduty.IncidentNote{
		ID: incidentID,
		User: pagerduty.APIObject{
			ID:      c.defaultUser.ID,
			Type:    typeUserReference,
			Summary: c.defaultUser.Email, //as we use api key which is not bound to a user, we need to give the email and not c.defaultUser.Summary,
			Self:    c.defaultUser.Self,
			HTMLURL: c.defaultUser.HTMLURL,
		},
		Content:   fmt.Sprintf("Incident was acknowledged on behalf of %s. time: %s", actualAcknowledger, now.String()),
		CreatedAt: now.String(),
	}

	return c.pagerdutyClient.CreateIncidentNoteWithResponse(incidentID, note)
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
		if util.NormalizeString(sched.Name) == util.NormalizeString(scheduleName) {
			return &sched, nil
		}
	}

	return nil, fmt.Errorf("schedule not found")
}
