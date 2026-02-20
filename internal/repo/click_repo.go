package repo

import (
	"fmt"
	"time"

	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

func (r *ClickRepo) clickFilters(q *gorm.DB, trackerID, campaignID, channelID string) *gorm.DB {
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	if campaignID != "" {
		q = q.Where("campaign_id = ?", campaignID)
	}
	if channelID != "" {
		q = q.Where("channel_id = ?", channelID)
	}
	return q
}

type ClickRepo struct {
	DB *gorm.DB
}

func NewClickRepo(db *gorm.DB) *ClickRepo {
	return &ClickRepo{DB: db}
}

func (r *ClickRepo) Create(c *models.Click) error {
	return r.DB.Create(c).Error
}

func (r *ClickRepo) GetByID(id string) (*models.Click, error) {
	var c models.Click
	err := r.DB.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClickRepo) CountByDay(start, end time.Time, trackerID, campaignID, channelID string) ([]DailyCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("DATE(ts) AS date, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	if campaignID != "" {
		q = q.Where("campaign_id = ?", campaignID)
	}
	if channelID != "" {
		q = q.Where("channel_id = ?", channelID)
	}
	var results []DailyCount
	err := q.Group("DATE(ts)").Order("date").Find(&results).Error
	return results, err
}

func (r *ClickRepo) Summary(start, end time.Time, trackerID, campaignID, channelID string) (total, uniqueVisitors, bots int64, err error) {
	q := r.DB.Model(&models.Click{}).
		Select("COUNT(*) AS total, COUNT(DISTINCT visitor_id) AS unique_visitors, SUM(CASE WHEN is_bot THEN 1 ELSE 0 END) AS bots").
		Where("ts BETWEEN ? AND ?", start, end)
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	if campaignID != "" {
		q = q.Where("campaign_id = ?", campaignID)
	}
	if channelID != "" {
		q = q.Where("channel_id = ?", channelID)
	}
	var row struct {
		Total          int64
		UniqueVisitors int64
		Bots           int64
	}
	err = q.Row().Scan(&row.Total, &row.UniqueVisitors, &row.Bots)
	return row.Total, row.UniqueVisitors, row.Bots, err
}

func (r *ClickRepo) TopByGroup(start, end time.Time, dimension string, limit int) ([]GroupCount, error) {
	var selectClause, joinClause, groupClause string
	switch dimension {
	case "tracker_id":
		selectClause = "clicks.tracker_id AS group_id, trackers.name AS name, COUNT(*) AS count"
		joinClause = "JOIN trackers ON trackers.id = clicks.tracker_id"
		groupClause = "clicks.tracker_id, trackers.name"
	case "campaign_id":
		selectClause = "clicks.campaign_id AS group_id, campaigns.name AS name, COUNT(*) AS count"
		joinClause = "JOIN campaigns ON campaigns.id = clicks.campaign_id"
		groupClause = "clicks.campaign_id, campaigns.name"
	case "channel_id":
		selectClause = "clicks.channel_id AS group_id, channels.name AS name, COUNT(*) AS count"
		joinClause = "JOIN channels ON channels.id = clicks.channel_id"
		groupClause = "clicks.channel_id, channels.name"
	default:
		return nil, fmt.Errorf("unsupported dimension: %s", dimension)
	}
	var results []GroupCount
	err := r.DB.Table("clicks").
		Select(selectClause).
		Joins(joinClause).
		Where("clicks.ts BETWEEN ? AND ? AND clicks.is_bot = false", start, end).
		Group(groupClause).
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

func (r *ClickRepo) TopReferrers(start, end time.Time, trackerID, campaignID, channelID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("referer AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.clickFilters(q, trackerID, campaignID, channelID)
	var raw []NameCount
	err := q.Group("referer").Order("count DESC").Limit(500).Find(&raw).Error
	if err != nil {
		return nil, err
	}
	// Normalize hostnames and re-aggregate
	hostMap := make(map[string]int64)
	for _, r := range raw {
		host := NormalizeReferrerHost(r.Name)
		hostMap[host] += r.Count
	}
	return mapToSortedNameCounts(hostMap, limit), nil
}

func (r *ClickRepo) RawUACounts(start, end time.Time, trackerID, campaignID, channelID string) ([]NameCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("ua AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ?", start, end)
	q = r.clickFilters(q, trackerID, campaignID, channelID)
	var results []NameCount
	err := q.Group("ua").Order("count DESC").Limit(500).Find(&results).Error
	return results, err
}

func (r *ClickRepo) LanguageDistribution(start, end time.Time, trackerID, campaignID, channelID string, limit int) ([]NameCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("lang AS name, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.clickFilters(q, trackerID, campaignID, channelID)
	var results []NameCount
	err := q.Group("lang").Order("count DESC").Limit(limit).Find(&results).Error
	return results, err
}

func (r *ClickRepo) BotCountByDay(start, end time.Time, trackerID, campaignID, channelID string) ([]DailyCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("DATE(ts) AS date, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = true", start, end)
	q = r.clickFilters(q, trackerID, campaignID, channelID)
	var results []DailyCount
	err := q.Group("DATE(ts)").Order("date").Find(&results).Error
	return results, err
}

func (r *ClickRepo) CountByHour(start, end time.Time, trackerID, campaignID, channelID string) ([]HourlyCount, error) {
	q := r.DB.Model(&models.Click{}).
		Select("EXTRACT(HOUR FROM ts)::int AS hour, COUNT(*) AS count").
		Where("ts BETWEEN ? AND ? AND is_bot = false", start, end)
	q = r.clickFilters(q, trackerID, campaignID, channelID)
	var results []HourlyCount
	err := q.Group("hour").Order("hour").Find(&results).Error
	return results, err
}
