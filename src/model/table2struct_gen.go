package model

import (
	"github.com/gohouse/converter"
)

// Public visibility
func AddStruct(tableName string) bool {
	// Initialization
	t2t := converter.NewTable2Struct()
	// Personalized configuration
	t2t.Config(&converter.T2tConfig{
		// If the initial of the field is capitalized, tag will not be added, false will be added by default, and true will not be added
		RmTagIfUcFirsted: false,
		// Whether the field name of tag is converted to lowercase? If it has uppercase letters, the default value is false
		TagToLower: false,
		// If you want to convert other letters to lowercase while the initial of the field is uppercase, false will not be converted by default
		UcFirstOnly: false,
		////Put each struct into a separate file, false by default, and put it into the same file (not provided yet)
		//SeperatFile: false,
	})
	// Start migration transformation
	err := t2t.
		// Specify a table. If not specified, all tables will be migrated by default
		Table(tableName).
		// Table prefix
		// Prefix("prefix_").
		// Add json tag or not
		EnableJsonTag(true).
		// Package name of the generated struct (if it is empty by default, it will be named: package model)
		PackageName("model").
		// The key value of the tag field, which is orm by default
		TagKey("orm").
		// Add structure method to get table name
		RealNameMethod(tableName).
		// Generated structure save path
		SavePath("src/model/model.go").
		// Database dsn, which can be replaced by t2t.DB(). The parameter is the * sql.DB object
		Dsn("root:@tcp(localhost:3306)/test?charset=utf8").
		// implement
		Run()

	if err != nil {
		return true
	}

	return false
}
