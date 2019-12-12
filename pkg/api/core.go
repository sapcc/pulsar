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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"github.com/sapcc/pulsar/pkg/auth"
	"github.com/sapcc/pulsar/pkg/config"
)

const (
	actionType            = "button"
	actionName            = "reaction"
	actionTypeAcknowledge = "acknowledge"
)

// API ...
type API struct {
	authorizer *auth.Authorizer
	cfg        *config.SlackConfig
	logger     log.Logger
}

// New returns a new API or an error.
func New(authorizer *auth.Authorizer, cfg *config.SlackConfig, logger log.Logger) (*API, error) {
	return &API{
		logger:     log.With(logger, "component", "api"),
		authorizer: authorizer,
		cfg:        cfg,
	}, nil
}

// Serve ...
func (a *API) Serve(stop <-chan struct{}) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", a.home)
	router.HandleFunc("/interaction", a.handleInteraction).Methods(http.MethodPost)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.APIHost, a.cfg.APIPort))
	if err != nil {
		level.Error(a.logger).Log("msg", "error creating listener", "err", err.Error())
		return
	}
	defer ln.Close()

	level.Info(a.logger).Log("msg", "serving API", "host", a.cfg.APIHost, "port", a.cfg.APIPort)
	go http.Serve(ln, router)
	<-stop
}

func (a *API) home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// handles interactive slack requests.
func (a *API) handleInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		level.Debug(a.logger).Log("msg", "invalid request with method", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		level.Error(a.logger).Log("msg", "error reading request body", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBody, err := url.QueryUnescape(string(buf))
	if err != nil {
		level.Error(a.logger).Log("msg", "error unescaping request body", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonBody), &message); err != nil {
		level.Error(a.logger).Log("msg", "error unmarshalling json body", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if message.Token != a.cfg.VerificationToken {
		level.Info(a.logger).Log("msg", "invalid verification token on message", "token", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !a.authorizer.IsUserAuthorized(message.User.ID) {
		level.Info(a.logger).Log("msg", "rejecting unauthorized user", "username", message.User.Name, "userid", message.User.ID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := handleActionCallback(message.ActionCallback); err != nil {
		level.Error(a.logger).Log("msg", "error handeling message", "err", err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

func handleActionCallback(actionCallbacks slack.ActionCallbacks) error {
	for _, act := range actionCallbacks.AttachmentActions {
		switch act.Name {
		case "":
			// TODO:
		}
	}

	for _, act := range actionCallbacks.BlockActions {
		//TODO:
		fmt.Println(act.Text)
	}

	return nil
}
