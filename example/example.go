package main

import (
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	fibertracing "github.com/ibenitez/fiber-opentracing"
	"github.com/ibenitez/fiber-opentracing/fjaeger"
	"github.com/opentracing/opentracing-go"
)

func main() {
	app := *fiber.New()

	fjaeger.New(fjaeger.Config{})

	app.Use(fibertracing.New(fibertracing.Config{
		Tracer: opentracing.GlobalTracer(),
		OperationName: func(ctx *fiber.Ctx) string {
			return "TEST:  HTTP " + ctx.Method() + " URL: " + ctx.Path()
		},
	}))

	go func() {
		interruptor := make(chan os.Signal, 1)
		signal.Notify(interruptor, os.Interrupt)
		for range interruptor {
			app.Shutdown()
			os.Exit(1)
		}
	}()

	api := app.Group("/api")
	api.Get("/", indexHandle)
	app.Listen(":8080")
}

func indexHandle(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"msg": "test"})
}
