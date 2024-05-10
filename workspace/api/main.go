package main

import (
	"fmt"
	"log"
	"os"

	"github.com/T1-SD/internal/database"
	"github.com/T1-SD/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	err := database.Init(os.Getenv("MONGODB_URI"), os.Getenv("DB_NAME"))

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB")

	defer func() {
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	r := router.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}
