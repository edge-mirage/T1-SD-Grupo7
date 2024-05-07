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

	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS
	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS
	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS
	err := database.Init(os.Getenv("MONGODB_URI"), "cluster0")
	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS
	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS
	//CAMBIAR EL NOMBRE DE LA BASE DE DATOS

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
