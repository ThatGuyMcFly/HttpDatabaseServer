package authentication

import (
	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationError struct{}

func (e AuthenticationError) Error() string {
	return "User Not Authenticated"
}

type ExpiredPasswordError struct{}

func (e ExpiredPasswordError) Error() string {
	return "Password is expired"
}

type AuthorizationError struct{}

func (e AuthorizationError) Error() string {
	return "User Not Authorized"
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func AuthenticateUser(employeeId int, password string) (string, error) {
	db := database.ConnectDatabase(database.EmployeeDatabase)

	storedPassword, expired := database.GetEmployeePassword(db, employeeId)
	if expired {
		return "", ExpiredPasswordError{}
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", AuthorizationError{}
	}

	authToken := uuid.New().String()

	return authToken, nil
}

func AuthorizeUser(employeeId int, authToken string) (bool, error) {

	return true, nil
}

func SetUserPassword(employeeId int, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	db := database.ConnectDatabase(database.EmployeeDatabase)

	err = database.AddEmployeePassword(db, employeeId, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
