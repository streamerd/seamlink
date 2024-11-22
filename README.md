# Seamlink - SEO-Friendly Link Tracking for Fiber

[![Go Reference](https://pkg.go.dev/badge/github.com/streamerd/seamlink.svg)](https://pkg.go.dev/github.com/streamerd/seamlink)
[![Go Report Card](https://goreportcard.com/badge/github.com/streamerd/seamlink)](https://goreportcard.com/report/github.com/streamerd/seamlink)

Seamlink is a lightweight, SEO-friendly link tracking middleware for [Fiber](https://github.com/gofiber/fiber) web framework. Track outbound clicks and referral sources without compromising SEO rankings or user experience.

## Features

1. 🔗 SEO-friendly outbound link tracking
2. 📊 Referrer source tracking
3. 🚀 Zero impact on page load performance
4. 🎯 Automatic script injection
5. 🛡️ Domain exclusion support
6. 📱 User agent tracking
7. ⏱️ Timestamp recording
8. 🔄 Non-blocking async tracking
9. 💫 Client-side navigation preservation

## Install

This middleware supports Fiber v2.

```bash
go get -u github.com/gofiber/fiber/v2
go get -u github.com/streamerd/seamlink
```

## Signature

```go
seamlink.New(config ...seamlink.SeamlinkConfig) fiber.Handler
```

## Config

| Property | Type | Description | Default |
|----------|------|-------------|----------|
| StoreLinkClick | `func(SeamlinkClick) error` | Callback function for outbound click events | `nil` |
| StorePageVisit | `func(PageVisit) error` | Callback function for page visit events | `nil` |
| ExcludeDomains | `[]string` | List of domains to exclude from tracking | `[]` |

### Data Types

```go
type SeamlinkClick struct {
    URL       string    `json:"url"`
    Referrer  string    `json:"referrer"`
    UserAgent string    `json:"userAgent"`
    Timestamp time.Time `json:"timestamp"`
}

type PageVisit struct {
    URL       string    `json:"url"`
    Referrer  string    `json:"referrer"`
    UserAgent string    `json:"userAgent"`
    Timestamp time.Time `json:"timestamp"`
}
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gofiber/fiber/v2"
    "github.com/streamerd/seamlink"
)

func main() {
    app := fiber.New()

    // Initialize Seamlink middleware
    app.Use(seamlink.New(seamlink.SeamlinkConfig{
        StoreLinkClick: func(click seamlink.SeamlinkClick) error {
            fmt.Printf("🔗 Outbound click: %s (from: %s)\n", click.URL, click.Referrer)
            return nil
        },
        StorePageVisit: func(visit seamlink.PageVisit) error {
            fmt.Printf("👋 New visit: Came from %s to %s\n", visit.Referrer, visit.URL)
            return nil
        },
        ExcludeDomains: []string{"internal.example.com"},
    }))

    // Serve your content
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Type("html").SendString(`
            <!DOCTYPE html>
            <html>
                <head><title>Seamlink Demo</title></head>
                <body>
                    <a href="https://github.com">GitHub</a>
                </body>
            </html>
        `)
    })

    log.Fatal(app.Listen(":3000"))
}
```

### Advanced Example with Database Integration

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/streamerd/seamlink"
    "gorm.io/gorm"
)

type ClickEvent struct {
    gorm.Model
    URL       string
    Referrer  string
    UserAgent string
}

func main() {
    app := fiber.New()
    db := initDatabase() // Your database initialization

    app.Use(seamlink.New(seamlink.SeamlinkConfig{
        StoreLinkClick: func(click seamlink.SeamlinkClick) error {
            return db.Create(&ClickEvent{
                URL:       click.URL,
                Referrer:  click.Referrer,
                UserAgent: click.UserAgent,
            }).Error
        },
        ExcludeDomains: []string{"internal.example.com"},
    }))

    // Your routes here...
}
```

## How It Works

1. **Script Injection**: Seamlink automatically injects a lightweight tracking script before the closing `</body>` tag
2. **Click Tracking**: Captures outbound link clicks using event listeners
3. **Source Tracking**: Records referrer information when users visit your site
4. **Async Processing**: Uses `fetch` API to send tracking data without blocking navigation
5. **SEO Preservation**: Maintains original link structures without redirects

## Best Practices

1. 🎯 Use domain exclusion for internal links
2. 📊 Implement proper error handling in storage callbacks
3. 🔒 Consider privacy regulations when storing user data
4. 💾 Use appropriate database indexes for high-traffic scenarios
5. 🚀 Keep storage callbacks lightweight to prevent blocking

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- [Fiber Web Framework](https://github.com/gofiber/fiber)
- Inspired by various analytics solutions while prioritizing SEO
