package clients

import (
	"testing"

	"github.com/sapcc/go-pagerduty"
	"github.com/stretchr/testify/assert"
)

func TestFilterIncidents(t *testing.T) {
	stimuli := []pagerduty.Incident{
		{APIObject: pagerduty.APIObject{Summary: summaryText}},
		{APIObject: pagerduty.APIObject{Summary: summaryTextWithLink}},
		{APIObject: pagerduty.APIObject{Summary: summaryTextMultiple}},
		{APIObject: pagerduty.APIObject{Summary: summaryTextMultipleNoDescription}},
	}

	expected := []pagerduty.Incident{
		{APIObject: pagerduty.APIObject{Summary: summaryTextMultiple}},
	}

	f := &Filter{
		Clusters:  []string{"eu-de-1"},
		Alertname: "VVOLDatastoreNotAccessibleFromHost",
	}

	res := f.FilterIncidents(stimuli)
	assert.EqualValues(t, expected, res, "the slices should have equal content")
}
