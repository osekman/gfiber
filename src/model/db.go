package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	env "github.com/joho/godotenv"
)

var db *sql.DB
var err error
var once sync.Once

// Map is a shortcut for map[string]interface{}, useful for JSON returns
// type Map map[string]interface{}

type Response struct {
	Status  int64  
	Message string
	//Data    []map[string]interface{}
	Data    []fiber.Map
}

type DriversDB struct {
	dbMysql     *sql.DB
	dbSqlServer *sql.DB
}

var instance *DriversDB

func Connect() *DriversDB {

	err := env.Load(".env")
	if err != nil {
		fmt.Printf("Some error occured. Err: %s", err)
	}

	// hostname := "localhost"
	// port := 3306
	// username := "root"
	// password := ""
	// database := "test"

	db_type 	:= os.Getenv("db_type")
	hostname 	:= os.Getenv("db_host")
	port 		:= os.Getenv("db_port")
	username 	:= os.Getenv("db_username")
	password 	:= os.Getenv("db_password")
	database 	:= os.Getenv("db_name")

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, hostname, port, database)
	mdb, err := sql.Open("mysql", connString)

	fmt.Println("Failed to connect", db_type)

	if err != nil {
		fmt.Println("Failed to connect", err)
	}

	// ms-sql db conn
	hostname = ""
	port = ""
	username = ""
	password = ""
	database = ""

	// Build connection string
	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", hostname, username, password, port, database)
	mdb2, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()
	err = mdb.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	instance = &DriversDB{dbMysql: mdb, dbSqlServer: mdb2}

	return instance
}

func init() {
	Connect()
}

func GetData(conn string, sqlString string, parameters []interface{}) Response {
	resp := Response{Status: 0, Message: "", Data: []fiber.Map{} }

	switch conn {
	case "mysql":
		db = instance.dbMysql
	case "sqlserver":
		db = instance.dbSqlServer
	}

	if err != nil {
		fmt.Println("ERRR", err)
		resp.Status = 500
		return resp
	}

	rows, err := db.Query(sqlString)
	if err != nil {
		resp.Status = 400
		return resp
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		resp.Status = 400
		return resp
	}
	count := len(columns)
	tableData := make([]fiber.Map, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(fiber.Map)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	resp.Status = 200
	resp.Message = "OK"
	resp.Data = tableData

	rows.Close()

	return resp
}

func SetData(conn string, sqlString string, parameters []interface{}) Response {

	resp := Response{Status: 0, Message: "", Data: []fiber.Map{} }

	switch conn {
	case "mysql":
		db = instance.dbMysql
	case "sqlserver":
		db = instance.dbSqlServer
	}

	if err != nil {
		fmt.Println("ERRR", err)
	}

	stmt, err := db.Query(sqlString, parameters...)
	if err != nil {
		resp.Status = 400
		return resp
	}

	resp.Status = 200
	resp.Message = "OK"

	stmt.Close()

	return resp

}
