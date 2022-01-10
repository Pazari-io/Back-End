package main

import (
	"crypto/tls"
	"log"

	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/routes"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// initalize everything
	database.InitDB()
	app := fiber.New()

	// set up routes
	routes.InitRoutes(app)

	if internal.GetKey("PRODUCTION") == "true" {

		// Certificate manager
		m := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// Replace with your domain
			HostPolicy: autocert.HostWhitelist("api.pazari.io"),
			// Folder to store the certificates
			Cache: autocert.DirCache("./certs"),
		}

		// TLS Config
		cfg := &tls.Config{
			// Get Certificate from Let's Encrypt
			GetCertificate: m.GetCertificate,
			NextProtos: []string{
				"http/1.1", "acme-tls/1",
			},
		}
		ln, err := tls.Listen("tcp", "0.0.0.0:443", cfg)
		if err != nil {
			panic(err)
		}

		log.Fatal(app.Listener(ln))

	}

	app.Listen("0.0.0.0:" + internal.GetKey("PORT"))
}
