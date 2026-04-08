package database

type DatabaseFileName int

const (
	EmployeeDatabase DatabaseFileName = iota
	InventoryDatabase
	SessionDatabase
)

func (dm DatabaseFileName) String() string {
	return []string{"employee.db", "inventory.db", "session.db"}[dm]
}
