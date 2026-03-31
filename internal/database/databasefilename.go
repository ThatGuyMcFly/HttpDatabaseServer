package database

type DatabaseFileName int

const (
	EmployeeDatabase DatabaseFileName = iota
	InventoryDatabase
)

func (dm DatabaseFileName) String() string {
	return []string{"employee.db", "inventory.db"}[dm]
}
