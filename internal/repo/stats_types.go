package repo

type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type GroupCount struct {
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
	Count   int64  `json:"count"`
}

type HourlyCount struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type NameCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
