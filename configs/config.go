package configs

import (
	"github.com/Dozie2001/Open-Move-Webhook/internal/db"
	// "github.com/Dozie2001/Open-Move-Webhook/internal/models"
	
	"fmt"

	"github.com/joho/godotenv"
)

func Load() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error: cannot find .env file in the project root")
	}

	// Setup database connection
	db.SetupDB()

	// // Initialize firebase
	// firebase.Initialize()


	// // Initialize mailgun
	// mailgun.Initialize()

	// // Initialize smtp
	// smtp.Initialize()

	// // Initialize mailersend
	// mailersend.Initialize()
}