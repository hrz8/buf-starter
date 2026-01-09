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

var employeeStatusesFromProto = map[altalunev1.EmployeeStatus]EmployeeStatus{
	altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE:   EmployeeStatusActive,
	altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE: EmployeeStatusInactive,
}

func EmployeeStatusFromProto(s altalunev1.EmployeeStatus) EmployeeStatus {
	if v, ok := employeeStatusesFromProto[s]; ok {
		return v
	}
	return EmployeeStatusActive
}

var employeeStatusesToProto = map[EmployeeStatus]altalunev1.EmployeeStatus{
	EmployeeStatusActive:   altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE,
	EmployeeStatusInactive: altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE,
}

func EmployeeStatusToProto(s EmployeeStatus) altalunev1.EmployeeStatus {
	if v, ok := employeeStatusesToProto[s]; ok {
		return v
	}
	return altalunev1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
}

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

func (m *Employee) ToEmployeeProto() *altalunev1.Employee {
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

type UpdateEmployeeInput struct {
	ProjectID  int64
	PublicID   string // Employee's public ID
	Name       string
	Email      string
	Role       string
	Department string
	Status     EmployeeStatus
}

type UpdateEmployeeResult struct {
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

type DeleteEmployeeInput struct {
	ProjectID int64
	PublicID  string // Employee's public ID
}
