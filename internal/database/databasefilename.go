package database

type FileName int

const (
	EmployeeDatabase FileName = iota
	InventoryDatabase
	SessionDatabase
)

func (dm FileName) String() string {
	return []string{"employee.db", "inventory.db", "session.db"}[dm]
}
