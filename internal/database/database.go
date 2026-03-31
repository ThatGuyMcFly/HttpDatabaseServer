package database

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

const driverName = "sqlite"
const databasePath = "assets/databases/"

type Database struct {
	DatabaseName DatabaseFileName
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

type Employee struct {
	EmployeeId int
	FirstName  string
	LastName   string
	Role       string
}

type Query struct {
	Key   string
	Value string
}

type AddEmployeeError struct{}

func (e AddEmployeeError) Error() string {
	return "Failed to add employee"
}

type InvalidRoleError struct{}

func (e InvalidRoleError) Error() string {
	return "Invalid role"
}

func ConnectDatabase(databaseName DatabaseFileName) Database {

	db, err := sql.Open(driverName, databasePath+databaseName.String())
	if err != nil {
		return Database{}
	}

	return Database{
		DatabaseName: databaseName,
		Database:     db,
	}
}

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

func constructUpdateQuery(tableName string, data string, queries []Query) string {
	return ""
}

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

func roleNameToRoleId(role string) int {
	lowerRole := strings.TrimSpace(strings.ToLower(role))

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

func AddEmployee(database Database, employeeData *json.Decoder) ([]int, error) {

	var newIds []int

	for {
		var temp Employee
		err := employeeData.Decode(&temp)

		if err != nil && err != io.EOF {
			log.Println(err)
			return nil, AddEmployeeError{}
		} else if err == io.EOF {
			break
		}

		roleId := roleNameToRoleId(temp.Role)

		if roleId == 0 {
			return nil, InvalidRoleError{}
		}

		var employeeMap = map[string]any{
			"firstName": temp.FirstName,
			"lastName":  temp.LastName,
			"roleId":    roleId,
		}

		insertQuery := constructInsertQuery("Employee", employeeMap)

		result, err := database.Database.Exec(insertQuery)
		if err != nil {
			log.Println(err)
			return nil, AddEmployeeError{}
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println(err)
		}

		log.Printf("Added Employee with ID: %d\n", id)

		newIds = append(newIds, int(id))
	}

	return newIds, nil
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

func GetEmployees(database Database, queryString string) []Employee {
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

	var employees []Employee
	for rows.Next() {
		var employee = Employee{}
		var roleNumber int
		err = rows.Scan(&employee.EmployeeId, &employee.FirstName, &employee.LastName, &roleNumber)
		if err != nil {
			return nil
		}

		employee.Role = getRoleTitle(database, roleNumber)
		employees = append(employees, employee)
	}

	return employees
}

func UpdateEmployee(database Database, employeeData *json.Decoder) error {
	return nil
}

func DeleteEmployee(database Database, queryString string) error {
	var queries = extractQueries(queryString)

	var deleteQuery = constructDeleteQuery("Employee", queries)

	_, err := database.Database.Exec(deleteQuery)

	if err != nil {
		return err
	}

	return nil
}

func GetItems(database Database, queries []Query) []Item {
	if database.Database == nil {
		return nil
	}

	if database.DatabaseName != InventoryDatabase {
		return nil
	}

	var items []Item

	return items
}
