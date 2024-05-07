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
	menu := `
Bienvenido al sistema de protección de archivos de DiSis.
Para utilizar la aplicación seleccione los números
correspondientes al menú.
	
===== Ingrese o Registrese =====
1) Ingreso
2) Registro
3) Salir
================================
`

	menu_loged := `
------ Menú principal ------
1) Clientes
2) Protección
3) Salir
----------------------------
`

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	for {
		option := MainMenu(menu)
		switch option {
		case 1:
			login(menu_loged)
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

func MainMenu(menu string) int {
	fmt.Println(menu)

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

	requestBody, err := json.Marshal(map[string]interface{}{
		"Name":      name,
		"Last_name": lastname,
		"Rut":       rut,
		"Email":     mail,
		"Password":  password,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("\n¡Registro exitoso!")
}

func login(menu_loged string) {
	fmt.Print("\nIngrese su correo: ")
	var mail string
	fmt.Scan(&mail)

	fmt.Print("Ingrese su contraseña: ")
	var password string
	fmt.Scan(&password)

	requestBody, err := json.Marshal(map[string]interface{}{
		"Email":    mail,
		"Password": password,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
		os.Exit(0)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		fmt.Println("\n¡Ingreso exitoso!")
		UserMainMenu(menu_loged)
	} else {
		fmt.Println("\nUsuario o Contraseña incorrecta")
		fmt.Println("Intentalo de nuevo!")
		return
	}
}

func UserMainMenu(menu_loged string) {
	for {
		fmt.Println(menu_loged)

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
