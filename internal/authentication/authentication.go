package authentication

import (
	"github.com/google/uuid"
)

type AuthenticationError struct{}

func (e AuthenticationError) Error() string {
	return "User Not Authenticated"
}

type AuthorizationError struct{}

func (e AuthorizationError) Error() string {
	return "User Not Authorized"
}

func AuthenticateUser(employeeId int, password string) (string, error) {
	authToken := uuid.New().String()

	return authToken, nil
}

func AuthorizeUser(employeeId int, authToken string) (bool, error) {

	return true, nil
}
