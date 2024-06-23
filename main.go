package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

type Response struct {
	Id          int    `json:"id"`
	ProductName string `json:"productName"`
}

func main() {

	var (
		host                = os.Getenv("hostname")
		port, _             = strconv.Atoi(os.Getenv("port"))
		user                = os.Getenv("username")
		password            = os.Getenv("password")
		dbname              = os.Getenv("dbName")
		response []Response = []Response{}
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	app := fiber.New()

	api := app.Group("/module/product", logger.New())

	api.Get("/check", func(c *fiber.Ctx) error {

		c.Set("Content-Type", fiber.MIMEApplicationJSON)

		response = []Response{}

		rows, err := db.Query("SELECT * FROM public.item ORDER BY id ASC LIMIT 100")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				Id          int
				ProductName string
			)
			if err := rows.Scan(&Id, &ProductName); err != nil {
				log.Fatal(err)
			}
			response = append(response, Response{
				Id:          Id,
				ProductName: ProductName,
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		bytes, _ := json.Marshal(response)

		return c.Send(bytes)
	})

	log.Fatal(app.Listen(":3000"))

}
