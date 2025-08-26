package employee

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// mapEmployeeToProto converts domain Employee to proto Employee
func mapEmployeeToProto(emp *Employee) *altalunev1.Employee {
	var status altalunev1.EmployeeStatus
	switch emp.Status {
	case EmployeeStatusActive:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE
	case EmployeeStatusInactive:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE
	default:
		status = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
	}

	return &altalunev1.Employee{
		Id:         emp.ID,
		Name:       emp.Name,
		Email:      emp.Email,
		Role:       emp.Role,
		Department: emp.Department,
		Status:     status,
		CreatedAt:  timestamppb.New(emp.CreatedAt),
		UpdatedAt:  timestamppb.New(emp.UpdatedAt),
	}
}

// mapEmployeesToProto converts slice of domain Employees to proto Employees
func mapEmployeesToProto(employees []*Employee) []*altalunev1.Employee {
	if employees == nil {
		return make([]*altalunev1.Employee, 0)
	}

	result := make([]*altalunev1.Employee, 0, len(employees))
	for _, emp := range employees {
		result = append(result, mapEmployeeToProto(emp))
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map roles
	if roles, ok := filters["roles"]; ok && roles != nil {
		result["roles"] = &altalunev1.FilterValues{Values: roles}
	} else {
		result["roles"] = &altalunev1.FilterValues{Values: []string{}}
	}

	// Map departments
	if departments, ok := filters["departments"]; ok && departments != nil {
		result["departments"] = &altalunev1.FilterValues{Values: departments}
	} else {
		result["departments"] = &altalunev1.FilterValues{Values: []string{}}
	}

	// Map statuses (always include as they're constants)
	if statuses, ok := filters["statuses"]; ok && statuses != nil {
		result["statuses"] = &altalunev1.FilterValues{Values: statuses}
	} else {
		result["statuses"] = &altalunev1.FilterValues{Values: []string{"active", "inactive"}}
	}

	return result
}
