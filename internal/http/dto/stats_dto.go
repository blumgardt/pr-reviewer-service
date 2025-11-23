package dto

type ReviewerStatsItem struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	AssignedCount int64  `json:"assigned_count"`
}

type ReviewerStatsResponse struct {
	Items []ReviewerStatsItem `json:"items"`
}
