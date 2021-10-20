package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB
var err error
var once sync.Once

type Response struct {
	Status  int64  
	Message string
	Data    []map[string]interface{}
}

type DriversDB struct {
	dbMysql     *sql.DB
	dbSqlServer *sql.DB
}

var instance *DriversDB

func Connect() *DriversDB {

	hostname := "localhost"
	port := 3306
	username := "root"
	password := ""
	database := "test"

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", username, password, hostname, port, database)
	mdb, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println("Failed to connect", err)
	}

	hostname = "176.236.208.126"
	port = 1433
	username = "sa"
	password = "@Z7N!@Xbm8!"
	database = "modaselvim"

	// Build connection string
	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", hostname, username, password, port, database)
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
	resp := Response{Status: 0, Message: "", Data: []map[string]interface{}{} }

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
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
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

	resp := Response{Status: 0, Message: "", Data: []map[string]interface{}{} }

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