package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main(){
	// Menu Program Logic
	// Menu output
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

	menu_clients := `
------------ Menú clientes ------------
1) Listar los clientes registrados
2) Obtener un cliente por ID
3) Obtener un cliente por RUT
4) Registrar un nuevo cliente
5) Actualizar datos de un cliente
6) Borrar un cliente por ID
7) Volver
	`

	for {

		//Prints
		fmt.Println(menu)

		//Option extracction and validation
		var opcion int
		fmt.Print("Ingrese su opción: ")
		_, err := fmt.Scan(&opcion)
		if err != nil {
			fmt.Println("Error al leer la opción:", err)
			os.Exit(1)
		}

		var login_successfull int

		//Case statement
		switch opcion {
		case 1:
			fmt.Println("Seleccionaste Ingreso.")
			fmt.Scanln()
			login_successfull = login()
			if login_successfull == 1 {
				fmt.Println("Usuario o Contraseña incorrecta")
				fmt.Println("Intentalo de nuevo!")
			}else if login_successfull == 0{
				for{
					//Prints
					fmt.Println(menu_loged)
					var opcion_logged int
					fmt.Print("Ingrese su opción: ")
					_, err := fmt.Scan(&opcion_logged)
					if err != nil {
						fmt.Println("Error al leer la opción:", err)
						os.Exit(1)
					}
					
					switch opcion_logged {
					case 1:
						// fmt.Println("Seleccionaste Registro.")
						fmt.Scanln()
					case 2:
						// fmt.Println("Seleccionaste Registro.")
						fmt.Scanln()
					case 3:
						fmt.Println("Saliendo del programa.")
						os.Exit(0)
					default:
						fmt.Println("Opción no válida. Por favor seleccione una opción válida.")
					}
				}
			}
		case 2:
			fmt.Println("Seleccionaste Registro.")
			fmt.Scanln()
			registUser()
		case 3:
			fmt.Println("Saliendo del programa.")
			os.Exit(0)
		default:
			fmt.Println("Opción no válida. Por favor seleccione una opción válida.")
		}
	}
}

func login() int{
	fmt.Println("++++++++ Login ++++++++")
	//Ask for name
	var email string
	fmt.Print("Ingrese el correo de su cuenta: ")
	fmt.Scanln(&email)

	//Ask for lastname
	var passwd string
	fmt.Print("Ingrese su contraseña: ")
	fmt.Scanln(&passwd)

	//creating the query
	requestBody, err := json.Marshal(map[string]string{
		"email":email,
		"password":passwd,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
		os.Exit(0)
	}

	// making the POST
	resp, err := http.Post("http://127.0.0.1:5000/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	// read de status code
	if resp.Status == "200 OK" {
		fmt.Println("¡Ingreso exitoso!")
		return 0
	}else{
		return 1
	}
}

func registUser(){
	fmt.Println("++++++++ Registro ++++++++")
	//Ask for name
	var name string
	fmt.Print("Ingrese su nombre: ")
	fmt.Scanln(&name)

	//Ask for lastname
	var lastname string
	fmt.Print("Ingrese su apellido: ")
	fmt.Scanln(&lastname)

	//Ask for RUT
	var rut string
	fmt.Print("Ingrese su RUT: ")
	fmt.Scanln(&rut)

	//Ask for mail
	var email string
	fmt.Print("Ingrese su correo: ")
	fmt.Scanln(&email)

	//Ask for password
	var passwd string
	fmt.Print("Ingrese su contraseña: ")
	fmt.Scanln(&passwd)

	//creating the query
	requestBody, err := json.Marshal(map[string]string{
		"name":name,
		"last_name":lastname,
		"rut":rut,
		"email":email,
		"password":passwd,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
		return
	}

	// making the POST
	resp, err := http.Post("http://127.0.0.1:5000/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()

	
	// show the data sumbited
	fmt.Println("\nDatos ingresados:")
	fmt.Printf("Nombre: %s\n", name)
	fmt.Printf("Apellido: %s\n", lastname)
	fmt.Printf("RUT: %s\n", rut)
	fmt.Printf("Correo: %s\n", email)
	fmt.Printf("Contraseña: %s\n", passwd)
	// read de status code
	fmt.Println("Código de estado:", resp.Status)
}