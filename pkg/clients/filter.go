package clients

import (
	"github.com/sapcc/go-pagerduty"
	"github.com/sapcc/pulsar/pkg/util"
)

// Filter can be used to filter Pagerduty incidents.
type Filter struct {
	// Alertname is the alertname to filter for.
	Alertname,

	// Severity of alerts/instance to filter for.
	Severity,

	// Fingerprint of the alert(s) to filter for.
	Fingerprint string

	// Clusters to filter for
	Clusters []string

	// limit is the number of items per response. Use
	limit *uint
}

// ClusterFilterFromText takes a string potentially containing cluster names and creates the filter accordingly.
func (f *Filter) ClusterFilterFromText(theString string) error {
	clusters, err := util.ParseClusterFromString(theString)
	if err != nil {
		return err
	}
	f.Clusters = clusters
	return nil
}

// AlertnameFilterFromText takes a string potentially containing an alertname and creates the filter accordingly.
func (f *Filter) AlertnameFilterFromText(theString string) error {
	_, alertname, err := parseRegionAndAlertnameFromText(theString)
	if err != nil {
		return err
	}

	f.Alertname = util.NormalizeString(alertname)
	return nil
}

// FilterIncidents does what it says.
func (f *Filter) FilterIncidents(incidents []pagerduty.Incident) []pagerduty.Incident {
	res := make([]pagerduty.Incident, 0)

	for _, inc := range incidents {
		region, alertname, err := parseRegionAndAlertnameFromText(inc.Summary)
		if err != nil {
			continue
		}

		keep := true
		if f.Clusters != nil && !util.Contains(f.Clusters, region) {
			keep = false
		}

		if f.Alertname != "" && util.NormalizeString(f.Alertname) != alertname {
			keep = false
		}

		if keep {
			res = append(res, inc)
		}
	}

	return res
}

// SetLimit sets the limit of items per response.
func (f *Filter) SetLimit(limit uint) {
	f.limit = &limit
}

// GetLimit returns the current limit of items per response.
func (f *Filter) GetLimit() uint {
	if f.limit != nil {
		return *f.limit
	}
	return 100
}
