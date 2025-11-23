package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	PullRequestStatus string
	ReviewersID       []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}
