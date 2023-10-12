package webSocketModels

type Employee struct {
	ID             uint
	DepartmentID   uint
	ManualID       uint16
	LastName       string
	FirstName      string
	PatronymicName string
	Active         bool
	SortNumber     uint
	Position       string
	PhoneNumber    string
	Telegram       string
	Email          string
	Login          string
	Password       string
}
