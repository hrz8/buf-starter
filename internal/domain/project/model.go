package project

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EnvironmentStatus string

const (
	EnvironmentStatusLive    EnvironmentStatus = "live"
	EnvironmentStatusSandbox EnvironmentStatus = "sandbox"
)

var environmentStatusesFromString = map[string]EnvironmentStatus{
	"live":    EnvironmentStatusLive,
	"sandbox": EnvironmentStatusSandbox,
}

func EnvironmentStatusFromString(s string) EnvironmentStatus {
	if v, ok := environmentStatusesFromString[s]; ok {
		return v
	}
	return EnvironmentStatusSandbox
}

var environmentStatusesesToString = map[EnvironmentStatus]string{
	EnvironmentStatusLive:    "live",
	EnvironmentStatusSandbox: "sandbox",
}

func EnvironmentStatusToString(s EnvironmentStatus) string {
	if v, ok := environmentStatusesesToString[s]; ok {
		return v
	}
	return "sandbox"
}

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

func (r *ProjectQueryResult) ToProject() *Project {
	return &Project{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		Timezone:    r.Timezone,
		Environment: r.Environment,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
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

func (m *Project) ToProjectProto() *altalunev1.Project {
	return &altalunev1.Project{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Timezone:    m.Timezone,
		Environment: string(m.Environment),
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

type CreateProjectInput struct {
	Name        string
	Description string
	Timezone    string
	Environment EnvironmentStatus
}

type CreateProjectResult struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	Timezone    string
	Environment EnvironmentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *CreateProjectResult) ToProject() *Project {
	return &Project{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		Timezone:    r.Timezone,
		Environment: r.Environment,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type UpdateProjectInput struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	Timezone    string
}

type UpdateProjectResult struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	Timezone    string
	Environment EnvironmentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *UpdateProjectResult) ToProject() *Project {
	return &Project{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		Timezone:    r.Timezone,
		Environment: r.Environment,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type DeleteProjectInput struct {
	PublicID string
}
