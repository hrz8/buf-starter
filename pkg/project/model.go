package project

import (
	"errors"
	"time"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type EnvironmentStatus string

const (
	EnvironmentStatusLive    EnvironmentStatus = "live"
	EnvironmentStatusSandbox EnvironmentStatus = "sandbox"
)

type ProjectQueryResult struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	Timezone    string
	Environment EnvironmentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Project struct {
	ID          string
	Name        string
	Description string
	Timezone    string
	Environment EnvironmentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
