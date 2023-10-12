package webSocketModels

type DepartmentRequest struct {
	ID         uint
	Name       string
	Active     bool
	SortNumber uint
}

type DepartmentResponse DepartmentRequest
