package main

import (
	"context"
	"database/sql"
	"fmt"
	"mygfiber/src/model"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	env "github.com/joho/godotenv"
)


func main() {

	// app := fiber.New(fiber.Config{
    //     Views: html.New("./src/view", ".html"),
    // })


	app := fiber.New()

	app.Use(cors.New())

	// Default middleware config
	app.Use(func(c *fiber.Ctx) error {
		
		dt := time.Now()
		//Format YYYY-MM-DD
		log_path := "./"+dt.Format("2006-02-01")+".log"

		file, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}

		newLine := "[ " + c.IP() +" ] " +c.Method()+ " - "+c.OriginalURL()
		_, err = fmt.Fprintln(file, newLine)

		if err != nil {
			fmt.Println(err)
			defer file.Close()
		}

		// Go to next middleware:
		return c.Next()
	})

	//app.Use(compress.New()) // default
	app.Use(compress.New(compress.Config{ Level: compress.LevelBestSpeed, })) // en hizli yanit ver


	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Render("index", fiber.Map{
	// 	   "Title": "Hello, World!",
	// 	})
	// })


	api := app.Group("/api") // /api
	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
        c.Set("Version", "v1")
        return c.Next()
    })

	app.Static("/","./statics")

	// APP rooting
	app.Get("/test", func(c *fiber.Ctx) error {
		err := env.Load(".env")
		if err != nil {
			fmt.Printf("Some error occured. Err: %s", err)
		}

		val := os.Getenv("stack")
		fmt.Println(val)

		val = os.Getenv("port")
		fmt.Println(val)

		return c.SendString(val)
	})

	v1.Get("/sql", func(c *fiber.Ctx) error {

		// var db *sql.DB

		// fmt.Println(controller.Greet())
		// fmt.Println(controller.Greet2())
		// fmt.Println(model.AddStruct("customers"))

		// db, err := model.DB_SQLServer()
		// if err != nil {
		// 	log.Fatal("Error not connected ", err.Error())
		// }

		// ctx := context.Background()
		// err = db.PingContext(ctx)
		// if err != nil {
		// 	log.Fatal(err.Error())
		// }
		// fmt.Printf("Connected!\n")

		// // // Read employees
		// count, err := ReadEmployees(db)
		// if err != nil {
		// 	log.Fatal("Error reading Employees: ", err.Error())
		// }
		// fmt.Printf("Read %d row(s) successfully.\n", count)

		/**
		  const ile := beraber kullanılamaz.
		  Yanlış kullanım: const isim := “Ali”
		  Doğru kullanım: const isim = “Ali”

		  for a := 0; a < 10; a++{
		      fmt.Printf("value of a: %d\n", a)
		  }
		*/

		return c.SendString("Hello, World 👋!")
	})

	// Routes
	v1.Get("/j", getUsers).Post("/j", createUser).Put("/j", updateUser)
	//.Delete("j/:id", json)

	// GET /flights/LAX-SFO
	v1.Get("/flights/:from/:to", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	app.Listen(":3000")
}

// Handler
func getUsers(c *fiber.Ctx) error {

	// fmt.Println(model.GetData("mysql", "SELECT * FROM customers limit 5", []string{}))
	// fmt.Println(model.GetData("sqlserver", "SELECT TOP 5 id, tsnumber, kargo FROM ztKargo", []string{}))

	type SomeStruct struct {
		Name string
		Age  uint8
	}

	var mdata = []SomeStruct{
		SomeStruct{Name: "Grame1", Age: 20},
	}

	// dizilerde ekleme yaparken ekli halini diziye return etmek lazım.
	mdata = append(mdata, SomeStruct{Name: "Grame2", Age: 21})
	mdata = append(mdata, SomeStruct{Name: "Grame3", Age: 22})

	somevars := []interface{}{888888, "ooooooosssssmmmmm"}

	//re := model.SetData("mysql", "insert into customers (CustomerId,CustomerUsername) values (?,?)", somevars)
	re := model.GetData("sqlserver", "SELECT TOP 5 id, tsnumber, kargo FROM ztKargo", somevars)
	//re := model.GetData("mysql", "SELECT * FROM customers limit 5", somevars)

	return c.JSON(re)
}

// Handler
func createUser(c *fiber.Ctx) error {

	type ReqBody struct{
		Token string `json:"token" xml:"token" form:"token"`
		Username string `json:"username" xml:"username" form:"username"`
		OrderCode string `json:"order_code" xml:"order_code" form:"order_code"`
	}

	p := new(ReqBody)

	if  err:= c.BodyParser(p); err != nil {
		return err
	}

	fmt.Println( "token  >> " + p.Token )
	fmt.Println( "username >> "+ p.Username )
	fmt.Println( "ordercode >> "+ p.OrderCode )

	type SomeStruct struct {
		Name string `json:"name" xml:"name" form:"name"`
		Age  uint8  `json:"age" xml:"age" form:"age"`
	}

	var data = []SomeStruct{
		SomeStruct{Name: "Grame1", Age: 20},
	}

	// dizilerde ekleme yaparken ekli halini diziye return etmek lazım.
	data = append(data, SomeStruct{Name: "Grame2", Age: 21})
	data = append(data, SomeStruct{Name: "Grame3", Age: 22})

	//data = append(data, SomeStruct{Name: p.Name, Age: p.Age})

	return c.JSON(data)
}

// Handler
func updateUser(c *fiber.Ctx) error {

	type SomeStruct struct {
		Id   int    `json:"id" xml:"id" form:"id"`
		Name string `json:"name" xml:"name" form:"name"`
		Age  uint8  `json:"age" xml:"age" form:"age"`
	}

	p := new(SomeStruct)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	var data = []SomeStruct{}
	// dizilerde ekleme yaparken ekli halini diziye return etmek lazım.
	data = append(data, SomeStruct{Id: 1, Name: "Grame2", Age: 21})
	data = append(data, SomeStruct{Id: 2, Name: "Grame3", Age: 22})

	for i, s := range data {
		fmt.Println(i, s)
		if s.Id == p.Id {
			data[i].Name = p.Name
		}
	}

	return c.JSON(data)
}

// ReadEmployees reads all employee records
func ReadEmployees(db *sql.DB) (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("SELECT TOP 10 id, tsnumber, kargo FROM ztKargo")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var id int
		var tsnumber string
		var kargo string

		// Get values from row.
		err := rows.Scan(&id, &tsnumber, &kargo)
		if err != nil {
			return -1, err
		}

		fmt.Printf("ID: %d", id)
		count++
	}

	return count, nil
}

func getFile() *os.File {

	dt := time.Now()
	//Format YYYY-MM-DD
	log_path := "./"+dt.Format("2006-02-01")+".log"

	f, err := os.OpenFile(log_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	f.Close()

	return f

}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
	
    return !info.IsDir()
}

// func GetData(sql string, params []string) ([]string, error){

// 	fmt.Printf("'%s' ------>> %s\n",sql,params)
// 	var names []interface{}
// 	names = append(names, {name:"ali", yas:12})

// 	return "", nil
// }
