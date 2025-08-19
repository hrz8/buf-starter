package employee

import "time"

type EmployeeStatus string

const (
	EmployeeStatusActive   EmployeeStatus = "active"
	EmployeeStatusInactive EmployeeStatus = "inactive"
)

type EmployeeQueryResult struct {
	ID         int64
	PublicID   string
	Name       string
	Email      string
	Role       string
	Department string
	Status     EmployeeStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Employee struct {
	ID         string
	Name       string
	Email      string
	Role       string
	Department string
	Status     EmployeeStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
