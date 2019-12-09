package clients

import (
	"github.com/sapcc/go-pagerduty"
	"github.com/sapcc/pulsar/pkg/util"
	"regexp"
)

const clusterRegex = `([\w-]*\w{2}-\w{2}-\d)|admin|staging`

// IncidentFilter can be used to filter Pagerduty incidents.
type IncidentFilter struct {
	Alertname,
	Severity,
	Fingerprint string
	Clusters []string
}

// ClusterFilterFromText takes a string potentially containing cluster names and creates the filter accordingly.
func (f *IncidentFilter) ClusterFilterFromText(theString string) error {
	r, err := regexp.Compile(clusterRegex)
	if err != nil {
		return err
	}

	f.Clusters = r.FindAllString(theString, -1)
	f.normalizeClusters()
	return nil
}

// FilterIncidents does what it says.
func (f *IncidentFilter) FilterIncidents(incidents []pagerduty.Incident) []pagerduty.Incident {
	res := make([]pagerduty.Incident, 0)

	for _, inc := range incidents {
		region, alertname, err := parseRegionAndAlertnameFromPagerdutySummary(inc.Summary)
		if err != nil {
			continue
		}

		keep := true
		if f.Clusters != nil && !util.Contains(f.Clusters, region) {
			keep = false
		}

		if f.Alertname != "" && normalizeString(f.Alertname) != alertname {
			keep = false
		}

		if keep {
			res = append(res, inc)
		}
	}

	return res
}

func (f *IncidentFilter) normalizeClusters() {
	for idx, c := range f.Clusters {
		f.Clusters[idx] = normalizeString(c)
	}
}
