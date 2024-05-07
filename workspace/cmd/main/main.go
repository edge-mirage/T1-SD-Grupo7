package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Bienvenido al sistema de protección de archivos de DiSis. Para utilizar la aplicación seleccione los números correspondientes al menú.")

	for {
		option := MainMenu()
		switch option {
		case 1:
			login()
		case 2:
			register()
		case 3:
			fmt.Println("\n¡Vuelve pronto!")
			return
		default:
			fmt.Println("\nOpción no válida. Por favor, intenta de nuevo.")
		}
	}
}

func MainMenu() int {
	fmt.Println("\nIngrese o regístrese")
	fmt.Println("1) Ingreso")
	fmt.Println("2) Registro")
	fmt.Println("3) Salir")

	return ReadOption()
}

func register() {
	fmt.Print("Ingrese su nombre: ")
	var name string
	fmt.Scan(&name)

	fmt.Print("Ingrese su apellido: ")
	var lastname string
	fmt.Scan(&lastname)

	fmt.Print("Ingrese su RUT: ")
	var rut string
	fmt.Scan(&rut)

	fmt.Print("Ingrese su correo: ")
	var mail string
	fmt.Scan(&mail)

	fmt.Print("Ingrese su contraseña: ")
	var password string
	fmt.Scan(&password)

	body := map[string]interface{}{
		"Name":      name,
		"Last_name": lastname,
		"Rut":       rut,
		"Email":     mail,
		"Password":  password,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/register", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("\n¡Registro exitoso!")
}

func login() {
	fmt.Print("\nIngrese su correo: ")
	var mail string
	fmt.Scan(&mail)

	fmt.Print("Ingrese su contraseña: ")
	var password string
	fmt.Scan(&password)

	fmt.Println("\n¡Ingreso exitoso!")

	body := map[string]interface{}{
		"Email":    mail,
		"Password": password,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	UserMainMenu()
}

func UserMainMenu() {
	for {
		fmt.Println("\nMenú principal")
		fmt.Println("1) Clientes")
		fmt.Println("2) Protección")
		fmt.Println("3) Salir")

		option := ReadOption()
		switch option {
		case 1:
			ClientMainMenu()
		case 2:
			protectFile()
		case 3:
			fmt.Println("Regresando al menú principal...")
			return
		default:
			fmt.Println("\nOpción no válida. Intenta de nuevo.")
		}
	}
}

func ReadOption() int {
	fmt.Print("\n")
	var option int
	fmt.Scan(&option)
	return option
}
