package seamlink

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func New(config ...SeamlinkConfig) fiber.Handler {
	fmt.Println("Seamlink middleware initialized")

	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	trackingScript := `
	<script>
	console.log('Seamlink tracking script loaded');
	
	const seamlinkTrack = function(e) {
		e.preventDefault();  // Prevent immediate navigation
		const link = e.currentTarget;
		const url = link.getAttribute('href');
		
		console.log('Link clicked:', url);
		
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
			window.open(url, '_blank');  // Open link after tracking
		}).catch(error => {
			console.error('Tracking error:', error);
			window.open(url, '_blank');  // Open link anyway if tracking fails
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

		// Handle tracking endpoint
		if c.Path() == "/api/seamlink/track" {
			fmt.Println("Received tracking request")

			var click SeamlinkClick
			if err := json.Unmarshal(c.Body(), &click); err != nil {
				fmt.Printf("Error parsing tracking data: %v\n", err)
				return c.Status(400).SendString("Invalid tracking data")
			}

			fmt.Printf("Processing click: %+v\n", click)

			if err := cfg.StoreLinkClick(click); err != nil {
				fmt.Printf("Error storing click: %v\n", err)
				return c.Status(500).SendString("Failed to store tracking data")
			}

			fmt.Println("Successfully processed click")
			return c.SendStatus(200)
		}

		// Process the response first
		if err := c.Next(); err != nil {
			return err
		}

		// Only proceed for HTML responses
		contentType := string(c.Response().Header.ContentType())
		if !strings.Contains(contentType, "text/html") {
			return nil
		}

		// Get the response body
		body := string(c.Response().Body())

		// Inject the script right before </body>
		if strings.Contains(body, "</body>") {
			fmt.Println("Injecting tracking script")
			modified := strings.Replace(body, "</body>", trackingScript+"</body>", 1)
			c.Response().SetBody([]byte(modified))
		} else {
			fmt.Println("No </body> tag found in response")
		}

		return nil
	}
}
