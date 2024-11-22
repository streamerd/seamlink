package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/streamerd/seamlink"
)

func main() {
	app := fiber.New()

	app.Use(seamlink.New(seamlink.SeamlinkConfig{
		StoreLinkClick: func(click seamlink.SeamlinkClick) error {
			fmt.Printf("💫 CLICK DETECTED: %s\n", click.URL)
			return nil
		},
		ExcludeDomains: []string{"stateful.art"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Serving index page")
		return c.Type("html").Send([]byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Seamlink Demo</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; }
        a { color: #0066cc; text-decoration: none; padding: 5px 10px; }
        a:hover { background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Seamlink Demo</h1>
    <p>Click these links to test outbound tracking (check server console):</p>
    <ul>
        <li><a href="https://github.com">GitHub</a></li>
        <li><a href="https://google.com">Google</a></li>
        <li><a href="https://stateful.art">start (excluded)</a></li>
        <li><a href="/about">Internal Page</a></li>
    </ul>
</body>
</html>`))
	})

	fmt.Println("🚀 Server starting on http://localhost:3000")
	fmt.Println("👀 Click tracking is active")

	log.Fatal(app.Listen(":3000"))
}
