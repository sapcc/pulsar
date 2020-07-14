Pulsar
-------------

[![Build Status](https://travis-ci.org/sapcc/pulsar.svg?branch=master)](https://travis-ci.org/sapcc/pulsar)

Interactive Slack bot for [Supernova](https://github.com/sapcc/supernova).

## Features

* List Prometheus alerts and Pagerduty incidents
* Acknowledge Pagerduty incidents
* List current Pagerduty on-call staff
* List Kubernetes nodes in a cluster 

## Installation

The [Pulsar Helm chart](https://github.com/sapcc/helm-charts/tree/master/global/pulsar) can be used for installation.

## Usage

The following secrets are provided via environment variables and are obtained after creating the Bot & enabling interactive messages in Slack via [this page](https://api.slack.com/apps).

```yaml
export SLACK_BOT_TOKEN = "topSecret!"
export SLACK_BOT_ID = "supernova"                                     
export SLACK_ACCESS_TOKEN = "superSecret?"
export SLACK_VERIFICATION_TOKEN = "anotherSecret!"
export SLACK_AUTHORIZED_USER_GROUP_NAMES = "slackGroup1,slackGroup2"
export SLACK_KUBERNETES_USER_GROUP_NAMES = "slackGroup3"
export SLACK_KUBERNETES_ADMIN_GROUP_NAMES = "slackGroup4"
export PAGERDUTY_DEFAULT_EMAIL = "defaultUser@pagerduty.com"
export PAGERDUTY_AUTH_TOKEN = "superSecret!"
```

## Development

Commands are independent plugins loaded during start and can be found in the [slack package](./pkg/slack).
See the [example command](./pkg/slack/hello.go).

### Release

Increment the [version](VERSION) and run `make release`.
