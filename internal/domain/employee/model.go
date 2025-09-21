package employee

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func (r *EmployeeQueryResult) ToEmployee() *Employee {
	return &Employee{
		ID:         r.PublicID,
		Name:       r.Name,
		Email:      r.Email,
		Role:       r.Role,
		Department: r.Department,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
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

func (m *Employee) ToEmployeeToProto() *altalunev1.Employee {
	var status altalunev1.EmployeeStatus
	switch m.Status {
	case EmployeeStatusActive:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE
	case EmployeeStatusInactive:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE
	default:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
	}

	return &altalunev1.Employee{
		Id:         m.ID,
		Name:       m.Name,
		Email:      m.Email,
		Role:       m.Role,
		Department: m.Department,
		Status:     status,
		CreatedAt:  timestamppb.New(m.CreatedAt),
		UpdatedAt:  timestamppb.New(m.UpdatedAt),
	}
}

type CreateEmployeeInput struct {
	ProjectID  int64
	Name       string
	Email      string
	Role       string
	Department string
	Status     EmployeeStatus
}

type CreateEmployeeResult struct {
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
