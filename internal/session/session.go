package session

import (
	"time"
)

type Session struct {
	Id           string
	EmployeeId   int
	AuthToken    string
	Created      time.Time
	LastAccessed time.Time
}

func StartSession(employeeId int, sessionToken string) {

}
