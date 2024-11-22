package seamlink

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// PageVisit tracks initial page visits
type PageVisit struct {
	URL       string    `json:"url"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
}

func New(config ...SeamlinkConfig) fiber.Handler {
	fmt.Println("Seamlink middleware initialized")

	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	trackingScript := `
	<script>
	console.log('Seamlink tracking script loaded');
	
	// Track page visit when loaded
	fetch('/api/seamlink/pageview', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			url: window.location.href,
			referrer: document.referrer || 'direct',
			userAgent: navigator.userAgent,
			timestamp: new Date().toISOString()
		})
	});
	
	// Track outbound clicks
	const seamlinkTrack = function(e) {
		e.preventDefault();
		const link = e.currentTarget;
		const url = link.getAttribute('href');
		
		fetch('/api/seamlink/track', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				url: url,
				referrer: window.location.href,
				userAgent: navigator.userAgent,
				timestamp: new Date().toISOString()
			})
		}).then(() => {
			window.open(url, '_blank');
		}).catch(error => {
			console.error('Tracking error:', error);
			window.open(url, '_blank');
		});
	};

	// Add tracking to links when DOM is ready
	document.addEventListener('DOMContentLoaded', function() {
		const links = document.querySelectorAll('a[href^="http"]');
		links.forEach(link => {
			link.addEventListener('click', seamlinkTrack);
			console.log('Added tracking to:', link.href);
		});
		console.log('Seamlink initialized on', links.length, 'links');
	});
	</script>`

	return func(c *fiber.Ctx) error {
		fmt.Println("Middleware called for path:", c.Path())

		// Handle page visit tracking
		if c.Path() == "/api/seamlink/pageview" {
			fmt.Println("Received pageview request")

			var visit PageVisit
			if err := json.Unmarshal(c.Body(), &visit); err != nil {
				fmt.Printf("Error parsing pageview data: %v\n", err)
				return c.Status(400).SendString("Invalid pageview data")
			}

			fmt.Printf("Processing pageview: %+v\n", visit)

			if err := cfg.StorePageVisit(visit); err != nil {
				fmt.Printf("Error storing pageview: %v\n", err)
				return c.Status(500).SendString("Failed to store pageview data")
			}

			fmt.Println("Successfully processed pageview")
			return c.SendStatus(200)
		}

		// Handle click tracking
		if c.Path() == "/api/seamlink/track" {
			fmt.Println("Received click tracking request")

			var click SeamlinkClick
			if err := json.Unmarshal(c.Body(), &click); err != nil {
				fmt.Printf("Error parsing click data: %v\n", err)
				return c.Status(400).SendString("Invalid click data")
			}

			fmt.Printf("Processing click: %+v\n", click)

			if err := cfg.StoreLinkClick(click); err != nil {
				fmt.Printf("Error storing click: %v\n", err)
				return c.Status(500).SendString("Failed to store click data")
			}

			fmt.Println("Successfully processed click")
			return c.SendStatus(200)
		}

		// Process the response
		if err := c.Next(); err != nil {
			return err
		}

		// Only proceed for HTML responses
		contentType := string(c.Response().Header.ContentType())
		if !strings.Contains(contentType, "text/html") {
			return nil
		}

		// Inject tracking script
		body := string(c.Response().Body())
		if strings.Contains(body, "</body>") {
			fmt.Println("Injecting tracking script")
			modified := strings.Replace(body, "</body>", trackingScript+"</body>", 1)
			c.Response().SetBody([]byte(modified))
		}

		return nil
	}
}
