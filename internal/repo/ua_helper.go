package repo

import (
	"net/url"
	"sort"

	"github.com/mssola/useragent"
)

// ParseUADistribution parses raw UA strings into browser and OS distributions.
func ParseUADistribution(uaCounts []NameCount, limit int) (browsers, oses []NameCount) {
	browserMap := make(map[string]int64)
	osMap := make(map[string]int64)

	for _, uc := range uaCounts {
		ua := useragent.New(uc.Name)
		browserName, _ := ua.Browser()
		if browserName == "" {
			browserName = "Unknown"
		}
		osInfo := ua.OSInfo()
		osName := osInfo.Name
		if osName == "" {
			osName = "Unknown"
		}
		browserMap[browserName] += uc.Count
		osMap[osName] += uc.Count
	}

	browsers = mapToSortedNameCounts(browserMap, limit)
	oses = mapToSortedNameCounts(osMap, limit)
	return
}

// NormalizeReferrerHost extracts the hostname from a URL string.
func NormalizeReferrerHost(ref string) string {
	if ref == "" {
		return "(direct)"
	}
	u, err := url.Parse(ref)
	if err != nil || u.Host == "" {
		return ref
	}
	return u.Host
}

func mapToSortedNameCounts(m map[string]int64, limit int) []NameCount {
	result := make([]NameCount, 0, len(m))
	for name, count := range m {
		result = append(result, NameCount{Name: name, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}
