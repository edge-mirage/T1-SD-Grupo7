package main

import (
	"fmt"
	"os"

	"github.com/T1-SD-Grupo7/database"
	"github.com/T1-SD-Grupo7/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
		return
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		fmt.Println("Missing environment variable: MONGODB_URI")
		return
	}

	err := database.Init(MONGODB_URI, "cluster0")
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	fmt.Println("Connected to MongoDB")

	defer func() {
		err := database.Close()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
	}()

	r := router.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}
