package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello, World! 2")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		// fmt.Println(c.Method(), c.Path())
		return c.Status(200).JSON((fiber.Map{"msg": "Hello, World!"}))
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		// fmt.Println(c.Method(), c.Path())
		return c.Status(200).JSON(todos)
	})

	// Create a Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // {id:0,completed:false,body:""} 값이 생긴다고 보면 된다. memory address를 얻어서 todo 에 할당

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo) // append는 값을 줘야하기 때문에 append(todos, todo)가 아니라 *todo 를 전달

		// var x int = 5   // 0x0001
		// var p *int = &x // 0x0001

		// fmt.Println(p)  // address
		// fmt.Println(*p) // value

		// 201: resource has been created
		return c.Status(201).JSON(todo)
	})

	// Update a Todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(400).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Delete a Todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"message": "Todo deleted"})
			}
		}

		return c.Status(400).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
