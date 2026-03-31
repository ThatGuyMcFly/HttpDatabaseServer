package database

type EmployeeRole int

const (
	Salesfloor EmployeeRole = iota
	Warehouse
	Administrator
)

func (e EmployeeRole) String() string {
	return []string{"Salesfloor", "Warehouse", "Administrator"}[e]
}
