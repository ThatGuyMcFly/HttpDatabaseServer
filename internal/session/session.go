package session

import (
	"time"
)

type Session struct {
	id           string
	employeeId   int
	authToken    string
	created      time.Time
	lastAccessed time.Time
}

func (session Session) GetId() string {
	return session.id
}

func (session Session) EetEmployeeId() int {
	return session.employeeId
}

func (session Session) GetAuthToken() string {
	return session.authToken
}

func (session Session) GetCreated() time.Time {
	return session.created
}

func (session Session) GetLastAccessed() time.Time {
	return session.lastAccessed
}

func StartSession(employeeId int, sessionToken string) Session {
	currentTime := time.Now()

	return Session{
		id:           "",
		employeeId:   employeeId,
		authToken:    sessionToken,
		created:      currentTime,
		lastAccessed: currentTime,
	}
}
