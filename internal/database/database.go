package database

import (
	"database/sql"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/ThatGuyMcFly/HttpDatabaseServer/internal/session"
	_ "modernc.org/sqlite"
)

const driverName = "sqlite"
const databasePath = "assets/databases/"

var notAnInt = int(math.NaN())

type Database struct {
	DatabaseName FileName
	Database     *sql.DB
}

type Item struct {
	ItemId          int
	ItemName        string
	Upc             string
	Description     string
	Category        string
	Price           float64
	SalesFloorStock int
	WarehouseStock  int
}

type Query struct {
	Key   string
	Value string
}

type AddEmployeeError struct{}

func (e AddEmployeeError) Error() string {
	return "Failed to add employee"
}

type InvalidKeyError struct{}

func (e InvalidKeyError) Error() string {
	return "Invalid key"
}

type InvalidRoleError struct{}

func (e InvalidRoleError) Error() string {
	return "Invalid role"
}

type InvalidDatabaseError struct{}

func (e InvalidDatabaseError) Error() string {
	return "Invalid database"
}

func ConnectDatabase(databaseName FileName) Database {

	db, err := sql.Open(driverName, databasePath+databaseName.String())
	if err != nil {
		return Database{}
	}

	return Database{
		DatabaseName: databaseName,
		Database:     db,
	}
}

func extractQueries(query string) []Query {
	if query == "" {
		return nil
	}

	var queries []Query

	var queryStrings = strings.Split(query, "&")

	for _, queryString := range queryStrings {
		var key = strings.Split(queryString, "=")[0]
		var value = strings.Split(queryString, "=")[1]
		queries = append(queries, Query{Key: key, Value: value})
	}

	return queries
}

//----------- Query Functions ------------//

func constructInsertQuery(tableName string, dataMap map[string]any) string {
	var insertQuery = INSERT.String() + " INTO " + tableName + " ("
	var columnString = ""
	var valueString = ""

	for key, value := range dataMap {
		if columnString != "" {
			columnString += ","
		}
		columnString += key

		if valueString != "" {
			valueString += ","
		}
		stringVal, stringOk := value.(string)
		if stringOk {
			valueString += "'" + stringVal + "'"
			continue
		}

		intVal, intOk := value.(int)
		if intOk {
			valueString += strconv.Itoa(intVal)
		}
	}

	insertQuery = insertQuery + columnString + ") VALUES (" + valueString + ")"

	return insertQuery
}

func constructSelectQuery(tableNames []string, columnNames []string, queries []Query) string {
	var selectQuery = SELECT.String()

	var columnString = ""

	if len(columnNames) == 0 || columnNames == nil {
		columnString = "*"
	} else {
		for _, columnName := range columnNames {
			if columnString != "" {
				columnString += ", "
			}
			columnString += columnName
		}
	}

	selectQuery += " " + columnString

	var tableString = ""

	for _, tableName := range tableNames {
		if tableString != "" {
			tableString += " "
		}
		tableString += tableName
	}

	selectQuery += " FROM " + tableString

	if len(queries) == 0 {
		return selectQuery
	}

	var parameterString = ""
	for _, query := range queries {
		if parameterString != "" {
			parameterString += " AND "
		}
		parameterString += query.Key + " LIKE '" + query.Value + "'"
	}

	return selectQuery + " WHERE " + parameterString
}

//func constructUpdateQuery(tableName string, data string, queries []Query) string {
//	return ""
//}

func constructDeleteQuery(tableName string, queries []Query) string {
	if len(queries) == 0 {
		return ""
	}

	deleteQuery := DELETE.String() + " FROM " + tableName + " WHERE "

	parameterString := ""
	for _, query := range queries {
		if parameterString != "" {
			parameterString += ", "
		}
		parameterString += query.Key + " = '" + query.Value + "'"
	}

	return deleteQuery + parameterString
}

//------------ Employee Functions ------------//

func roleNameToRoleId(role any) int {

	roleString, ok := role.(string)

	if !ok {
		return 0
	}

	lowerRole := strings.TrimSpace(strings.ToLower(roleString))

	switch lowerRole {
	case "administrator":
		return 1
	case "admin":
		return 1
	case "warehouse":
		return 2
	case "salesfloor":
		return 3
	default:
		return 0
	}
}

func getRoleTitle(database Database, roleId int) string {
	var query = "Select roleTitle from Role where roleId = ?"

	var rows, err = database.Database.Query(query, roleId)
	if err != nil {
		return ""
	}
	defer rows.Close()

	var title string
	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			return ""
		}

		return title
	}

	return ""
}

func createEmployeeMap(employee map[string]any) (map[string]any, error) {
	var employeeMap = map[string]any{
		"firstName": "",
		"lastName":  "",
		"roleId":    0,
	}

	expectedKeys := []string{
		"employeeId",
		"firstName",
		"lastName",
		"role",
	}

	for _, expectedKey := range expectedKeys {
		value, ok := employee[expectedKey]
		if !ok {
			if expectedKey == "employeeId" {
				continue
			}

			return nil, InvalidKeyError{}
		}

		if expectedKey == "role" {
			value = roleNameToRoleId(value)

			if value == 0 {
				return nil, InvalidRoleError{}
			}

			expectedKey = "roleId"
		}

		employeeMap[expectedKey] = value
	}

	return employeeMap, nil
}

func AddEmployee(database Database, employee map[string]any) (int, error) {

	employeeMap, err := createEmployeeMap(employee)
	if err != nil {
		return 0, err
	}

	insertQuery := constructInsertQuery("Employee", employeeMap)

	result, err := database.Database.Exec(insertQuery)
	if err != nil {
		log.Println(err)
		return 0, AddEmployeeError{}
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	log.Printf("Added Employee with ID: %d\n", id)

	return int(id), nil
}

func GetEmployees(database Database, queryString string) []map[string]any {
	if database.Database == nil {
		return nil
	}

	if database.DatabaseName != EmployeeDatabase {
		return nil
	}

	var queries = extractQueries(queryString)

	var tableNames = []string{
		"Employee",
	}

	var columnNames = []string{
		"employeeId",
		"firstName",
		"lastName",
		"roleId",
	}

	var selectQuery = constructSelectQuery(tableNames, columnNames, queries)

	var rows, err = database.Database.Query(selectQuery)
	if err != nil {
		return nil
	}

	var employees []map[string]any
	for rows.Next() {
		var employeeId, roleNumber int
		var firstName, lastName string

		err = rows.Scan(&employeeId, &firstName, &lastName, &roleNumber)
		if err != nil {
			return nil
		}

		role := getRoleTitle(database, roleNumber)

		employee := map[string]any{
			"employeeId": employeeId,
			"firstName":  firstName,
			"lastName":   lastName,
			"role":       role,
		}

		employees = append(employees, employee)
	}

	return employees
}

//func UpdateEmployee(database Database, employeeData *json.Decoder) error {
//	return nil
//}

func DeleteEmployee(database Database, queryString string) error {
	var queries = extractQueries(queryString)

	var deleteQuery = constructDeleteQuery("Employee", queries)

	_, err := database.Database.Exec(deleteQuery)

	if err != nil {
		return err
	}

	return nil
}

func AddEmployeePassword(database Database, employeeId int, password string) error {
	var passwordMap = map[string]any{
		"employeeId": employeeId,
		"password":   password,
		"expired":    0,
	}

	insetQuery := constructInsertQuery("Password", passwordMap)

	_, err := database.Database.Exec(insetQuery)
	if err != nil {
		return err
	}

	return nil
}

func GetEmployeePassword(database Database, employeeId int) (string, bool) {
	var tableNames = []string{
		"Password",
	}

	var columns = []string{
		"password",
		"expired",
	}

	queries := extractQueries("employeeId=" + strconv.Itoa(employeeId))

	selectQuery := constructSelectQuery(tableNames, columns, queries)
	var rows, err = database.Database.Query(selectQuery)
	if err != nil {
		return "", false
	}
	var password = ""
	var expired = 0
	for rows.Next() {
		rows.Scan(&password, &expired)
		if expired != 0 {
			return "", true
		}

		return password, false
	}

	return "", false
}

//------------ Item Functions ------------//

func GetItems(database Database, queries []Query) ([]Item, error) {
	if database.Database == nil {
		return nil, InvalidDatabaseError{}
	}

	if database.DatabaseName != InventoryDatabase {
		return nil, InvalidDatabaseError{}
	}

	var items []Item

	return items, nil
}

//------------ Session Functions ------------//

func AddSession(database Database, session session.Session) error {
	if database.Database == nil {
		return InvalidDatabaseError{}
	}

	if database.DatabaseName != SessionDatabase {
		return InvalidDatabaseError{}
	}

	sessionData := map[string]any{
		"employeeId":      session.EetEmployeeId(),
		"authToken":       session.GetAuthToken(),
		"datetimeCreated": session.GetCreated(),
		"lastAccessed":    session.GetLastAccessed(),
	}

	addSessionQuery := constructInsertQuery("Session", sessionData)

	_, err := database.Database.Exec(addSessionQuery)

	if err != nil {
		return err
	}

	return nil
}

func GetSession(database Database, authToken string) ([]session.Session, error) {
	if database.Database == nil {
		return nil, InvalidDatabaseError{}
	}

	if database.DatabaseName != SessionDatabase {
		return nil, InvalidDatabaseError{}
	}

	sessionQuery := "Select * from Session where authToken = ?"

	rows, err := database.Database.Query(sessionQuery, authToken)
	if err != nil {
		return nil, err
	}

	var sessions []session.Session

	for rows.Next() {
		var currentSession = session.Session{}

		err = rows.Scan(&currentSession)
		if err != nil {
			continue
		}

		sessions = append(sessions, currentSession)
	}

	return sessions, nil
}
