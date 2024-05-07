package main

import (
	"fmt"
	"os"

	"github.com/T1-SD-Grupo7/database"
	"github.com/T1-SD-Grupo7/router"
	"github.com/joho/godotenv"
)

func main() {
	//Loading de .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
		return
	}

	//Loading the db
	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		fmt.Println("Missing environment variable: MONGODB_URI")
		return
	}

	//Handling errors in the db
	err := database.Init(MONGODB_URI, "cluster0")
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	fmt.Println("Connected to MongoDB")

	//Handling the error to the console output
	defer func() {
		err := database.Close()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
	}()

	//caling router in ./router/router.go
	r := router.SetupRouter()
	r.Run(":" + os.Getenv("PORT"))
}