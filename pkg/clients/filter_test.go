package clients

import (
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
)
const (
	summaryText                      = "[#1594] \n [EU-DE-1] OpenstackLbaasApiFlapping - lbaas API flapping\n"
	summaryTextWithLink              = "[#1598] \n [AP-SA-1] BaremetalIronicSensorCritical - Sensor Critical for instance node009r-bm020.cc.ap-sa-1.cloud.sap\n"
	summaryTextMultiple              = "[#2130] \n [7 Alerts] [EU-DE-2] VVOLDatastoreNotAccessibleFromHost - vVOL Datastore accessibility check from host\n"
	summaryTextMultipleNoDescription = "[#2144] \n [3 Alerts] [EU-NL-1] OpenstackNeutronDatapathDown - \n"
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
		Clusters:  []string{"eu-de-2"},
		Alertname: "VVOLDatastoreNotAccessibleFromHost",
	}

	res := f.FilterIncidents(stimuli)
	assert.EqualValues(t, expected, res, "the slices should have equal content")
}
