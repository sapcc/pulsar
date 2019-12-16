package slack

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestListNodesCommand(t *testing.T) {
	cmd := listNodesCommand{}

	err := cmd.Init()
	assert.NoError(t, err, "there should be no error initializing the command")

	input := slack.Msg{Text: "list nodes in cluster s-qa-de-1"}
	_, err = cmd.Run(input)
	assert.NoError(t, err, "there should be no error listing the nodes in cluster s-qa-de-1")
}
