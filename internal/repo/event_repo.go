package repo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

func (r *EventRepo) eventFilters(q *gorm.DB, siteID string) *gorm.DB {
	if siteID != "" {
		q = q.Where("site_id = ?", siteID)
	}
	return q
}

type EventRepo struct {
	DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{DB: db}
}

func (r *EventRepo) BatchCreate(events []models.Event) error {
	return r.DB.CreateInBatches(events, 100).Error
}

func (r *EventRepo) CountByDay(start, end time.Time, siteID string) ([]DailyCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("DATE(ts) AS date, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	if siteID != "" {
		q = q.Where("site_id = ?", siteID)
	}
	var results []DailyCount
	err := q.Group("DATE(ts)").Order("date").Find(&results).Error
	return results, err
}

func (r *EventRepo) Summary(start, end time.Time, siteID string) (total, uniqueVisitors, uniqueSessions, bots int64, err error) {
	q := r.DB.Model(&models.Event{}).
		Select("COUNT(*) AS total, COUNT(DISTINCT visitor_id) AS unique_visitors, COUNT(DISTINCT session_id) AS unique_sessions, SUM(CASE WHEN is_bot THEN 1 ELSE 0 END) AS bots").
		Where("ts BETWEEN ? AND ?", start, end)
	if siteID != "" {
		q = q.Where("site_id = ?", siteID)
	}
	var row struct {
		Total          int64
		UniqueVisitors int64
		UniqueSessions int64
		Bots           int64
	}
	err = q.Row().Scan(&row.Total, &row.UniqueVisitors, &row.UniqueSessions, &row.Bots)
	return row.Total, row.UniqueVisitors, row.UniqueSessions, row.Bots, err
}

func (r *EventRepo) TopByGroup(start, end time.Time, dimension string, limit int) ([]GroupCount, error) {
	var selectClause, joinClause, groupClause string
	switch dimension {
	case "site_id":
		selectClause = "events.site_id AS group_id, sites.name AS name, COUNT(*) AS count"
		joinClause = "JOIN sites ON sites.id = events.site_id"
		groupClause = "events.site_id, sites.name"
	case "type":
		selectClause = "events.type AS group_id, events.type AS name, COUNT(*) AS count"
		joinClause = ""
		groupClause = "events.type"
	default:
		return nil, fmt.Errorf("unsupported dimension: %s", dimension)
	}
	q := r.DB.Table("events").
		Select(selectClause).
		Where("events.ts BETWEEN ? AND ? AND events.is_bot = false", start, end)
	if joinClause != "" {
		q = q.Joins(joinClause)
	}
	var results []GroupCount
	err := q.Group(groupClause).
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

func (r *EventRepo) TopReferrers(start, end time.Time, siteID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("referrer AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.eventFilters(q, siteID)
	var raw []NameCount
	err := q.Group("referrer").Order("count DESC").Limit(500).Find(&raw).Error
	if err != nil {
		return nil, err
	}
	hostMap := make(map[string]int64)
	for _, r := range raw {
		host := NormalizeReferrerHost(r.Name)
		hostMap[host] += r.Count
	}
	return mapToSortedNameCounts(hostMap, limit), nil
}

func (r *EventRepo) TopPages(start, end time.Time, siteID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("url AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.eventFilters(q, siteID)
	var raw []NameCount
	err := q.Group("url").Order("count DESC").Limit(500).Find(&raw).Error
	if err != nil {
		return nil, err
	}
	// Strip protocol+host, keep path
	pathMap := make(map[string]int64)
	for _, r := range raw {
		path := r.Name
		if u, err := url.Parse(r.Name); err == nil && u.Path != "" {
			path = u.Path
		}
		pathMap[path] += r.Count
	}
	return mapToSortedNameCounts(pathMap, limit), nil
}

func (r *EventRepo) RawUACounts(start, end time.Time, siteID string) ([]NameCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("ua AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ?", start, end)
	q = r.eventFilters(q, siteID)
	var results []NameCount
	err := q.Group("ua").Order("count DESC").Limit(500).Find(&results).Error
	return results, err
}

func (r *EventRepo) LanguageDistribution(start, end time.Time, siteID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("lang AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.eventFilters(q, siteID)
	var results []NameCount
	err := q.Group("lang").Order("count DESC").Limit(limit).Find(&results).Error
	return results, err
}

func (r *EventRepo) CountryDistribution(start, end time.Time, siteID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("country AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false AND country != ''", start, end)
	q = r.eventFilters(q, siteID)
	var results []NameCount
	err := q.Group("country").Order("count DESC").Limit(limit).Find(&results).Error
	return results, err
}

func (r *EventRepo) BotCountByDay(start, end time.Time, siteID string) ([]DailyCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("DATE(ts) AS date, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = true", start, end)
	q = r.eventFilters(q, siteID)
	var results []DailyCount
	err := q.Group("DATE(ts)").Order("date").Find(&results).Error
	return results, err
}

func (r *EventRepo) CountByHour(start, end time.Time, siteID string) ([]HourlyCount, error) {
	q := r.DB.Model(&models.Event{}).
		Select("EXTRACT(HOUR FROM ts)::int AS hour, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.eventFilters(q, siteID)
	var results []HourlyCount
	err := q.Group("hour").Order("hour").Find(&results).Error
	return results, err
}
