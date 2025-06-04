package main

import (
	"log"
	"os"

	"github.com/chat-socio/backend/cmd/app"
	"github.com/chat-socio/backend/cmd/migrate"
	"github.com/chat-socio/backend/configuration"
	"github.com/spf13/cobra"
	_ "github.com/chat-socio/backend/docs" // This is important!
)

// @title Chat Socio API
// @version 1.0
// @description This is the API documentation for Chat Socio backend service
// @host localhost:8887
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var (
	svc        string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "socio",
	Short: "Socio CLI is a command line interface for Socio",
	Long:  `Socio CLI is a command line interface for Socio`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		err := configuration.LoadConfig(configPath)
		if err != nil {
			log.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if svc == "" {
			log.Println("Please specify a service to run")
			os.Exit(1)
		}

		switch svc {
		case "app":
			// Run the app service
			app.RunApp()
		case "migrate":
			// Run the migration service
			migrate.Migrate()
		default:
			log.Printf("Unknown service: %s\n", svc)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&svc, "service", "s", "", "Service to run (app, migrate)")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "./config.yaml", "Path to the config file")

	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
	log.Println("Socio CLI executed successfully")
	os.Exit(0)
}
