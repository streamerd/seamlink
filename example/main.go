package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/streamerd/seamlink"
)

func main() {

	engine := html.New("./views", ".html")
	engine.Reload(true)
	engine.Debug(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(seamlink.New(seamlink.SeamlinkConfig{
		StoreLinkClick: func(click seamlink.SeamlinkClick) error {
			fmt.Printf("🔗 OUTBOUND CLICK:\n")
			fmt.Printf("   URL: %s\n", click.URL)
			fmt.Printf("   From Page: %s\n", click.Referrer)
			fmt.Printf("   User Agent: %s\n", click.UserAgent)
			return nil
		},
		StorePageVisit: func(visit seamlink.PageVisit) error {
			fmt.Printf("👋 NEW VISIT:\n")
			fmt.Printf("   To: %s\n", visit.URL)
			fmt.Printf("   From: %s\n", visit.Referrer)
			fmt.Printf("   User Agent: %s\n", visit.UserAgent)
			return nil
		},
		ExcludeDomains: []string{"internal.example.com"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Type("html").Send([]byte(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Tracked Page</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 40px auto;
            padding: 20px;
            background: #f7f7f7;
        }
        .content {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        a {
            display: inline-block;
            padding: 10px 20px;
            background: #2196F3;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            transition: background 0.3s;
        }
        a:hover {
            background: #1976D2;
        }
        .info {
            margin-top: 20px;
            padding: 15px;
            background: #e3f2fd;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="content">
        <h1>Tracked Page</h1>
        <p>This page tracks where you came from and where you go.</p>
        <a href="https://github.com">Visit GitHub</a>
        <div class="info">
            <p><strong>Check your terminal to see:</strong></p>
            <ul>
                <li>The referrer that brought you here</li>
                <li>Tracking when you click the GitHub link</li>
            </ul>
        </div>
    </div>
</body>
</html>`))
	})

	app.Get("/ref", func(c *fiber.Ctx) error {
		return c.Render("referrer", nil)
	})

	fmt.Println("🚀 Server starting on http://localhost:3000")
	fmt.Println("📊 Tracking both inbound and outbound traffic")

	log.Fatal(app.Listen(":3000"))
}
