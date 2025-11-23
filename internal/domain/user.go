package domain

type UserStatus string

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}
