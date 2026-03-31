package database

type SQLCommands int

const (
	SELECT SQLCommands = iota
	INSERT
	UPDATE
	DELETE
)

func (s SQLCommands) String() string {
	return []string{"SELECT", "INSERT", "UPDATE", "DELETE"}[s]
}
