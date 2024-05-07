package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
)

func main(){
	// Menu Program Logic
	// Menu output
	menuPrint := `
Bienvenido al sistema de protección de archivos de DiSis.
Para utilizar la aplicación seleccione los números
correspondientes al menú.

===== Ingrese o Registrese =====
1) Ingreso
2) Registro
3) Salir
================================
`

	//menu loop
	for {

		//Prints menu options
		fmt.Println(menuPrint)

		//Option extracction and validation
		opcion := ask4("Ingrese su opción")
		
		//Case statement
		switch opcion {
		case "1":
			menuLogin()
		case "2":
			menuRegister()
		case "3":
			fmt.Println("Saliendo del programa.")
			os.Exit(0)
		default:
			fmt.Println("Opción no válida. Por favor seleccione una opción válida.")
		}
	}
}

func ask4(mensaje string) string {
	var dato string
	fmt.Print(mensaje + ": ")
	_, err := fmt.Scanln(&dato)
	if err != nil {
		fmt.Println("Error al leer la opción:", err)
		os.Exit(1)
	}
	return dato
}

func menuLogin() {
	menuPrint := `
------ Menú principal ------
1) Clientes
2) Protección
3) Volver
----------------------------
`

	fmt.Println("Seleccionaste Ingreso.")
	fmt.Println("++++++++ Login ++++++++")

	email := ask4("Ingrese el correo")
	passwd := ask4("Ingrese su contraseña")
	
	//handling the login
	if 1 == login(email, passwd) { //1 mean wrong auth
		fmt.Println("Usuario o Contraseña incorrecta")
		fmt.Println("Intentalo de nuevo!")
	}else { //0 mean successfull otherwise exit the program with exit code 1
		
		for{// loop for the second menu
			
			//Print the menu for the loged person
			fmt.Println(menuPrint)
			
			//ask for the option
			opcion := ask4("Ingrese su opción")

			switch opcion{
			case "1":
				clientsMenu()
			case "2":
				//idClient := ask4("Escriba el ID del cliente objetivo")
				//path := ask4("Escriba la ruta donde se encuentra el archivo (incluya el nombre)")
			case "3":
				return
			default:
				fmt.Println("Opción no válida. Por favor seleccione una opción válida.")
			}
		}
	}
}

func menuRegister(){
	fmt.Println("Seleccionaste Registro.")
	fmt.Println("++++++++ Registro ++++++++")

	name := ask4("Ingrese su nombre")
	lastname := ask4("Ingrese su apellido")
	rut := ask4("Ingrese su RUT")
	email := ask4("Ingrese el correo")
	passwd := ask4("Ingrese su contraseña")

	if 1 == registUser(name, lastname, rut, email, passwd){ //1 mean wrong auth
		fmt.Println("El usuario ya existe")
		fmt.Println("Intenta Iniciar sesion")
	}
}

func clientsMenu(){
	menuPrint := `
------------ Menú clientes ------------
1) Listar los clientes registrados
2) Obtener un cliente por ID
3) Obtener un cliente por RUT
4) Registrar un nuevo cliente
5) Actualizar datos de un cliente
6) Borrar un cliente por ID
7) Volver
	`

	fmt.Println(menuPrint)
	//ask for the option
	opcion := ask4("Ingrese su opción")

	switch opcion{
	case "1"://list client
		resp, err := http.Get("http://127.0.0.1:5000/api/clients")
		if err != nil {
			fmt.Println("Error al hacer la solicitud GET:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Leer la respuesta
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer la respuesta:", err)
			os.Exit(1)
		}

		// Decodificar el JSON
		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error al decodificar el JSON:", err)
			os.Exit(1)
		}

		// Imprimir los campos de cada JSON
		for _, item := range data {
			fmt.Println("---")
			for key, value := range item {
				fmt.Printf("%s: %v\n", key, value)
			}
			fmt.Println("---")
		}
	case "2":
		//get client by id
	case "3":
		//get client by rut
	case "4":	
		// register new client
	case "5":
		//update client
	case "6":
		//delete client	
	case "7":
		//return to the previous menu
		fmt.Println("Devolviendose...")
		return
	default:
		fmt.Println("Opción no válida. Por favor seleccione una opción válida.")
	}

}

func login(email string, passwd string) int{
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

func registUser(name string, lastname string, rut string, email string, passwd string) int{
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
		os.Exit(0)
	}

	// making the POST
	resp, err := http.Post("http://127.0.0.1:5000/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	
	// show the data sumbited
	fmt.Println("\nDatos ingresados:")
	fmt.Printf("Nombre: %s\n", name)
	fmt.Printf("Apellido: %s\n", lastname)
	fmt.Printf("RUT: %s\n", rut)
	fmt.Printf("Correo: %s\n", email)
	fmt.Printf("Contraseña: %s\n", passwd)
	
	// return code
	if resp.Status == "201 Created" {
		fmt.Println("Registro exitoso!")
		return 0
	}else{
		return 1
	}
}

func registClient(name string, lastname string, rut string, email string) int{
	//creating the query
	requestBody, err := json.Marshal(map[string]string{
		"name":name,
		"last_name":lastname,
		"rut":rut,
		"email":email,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
		os.Exit(0)
	}

	// making the POST
	resp, err := http.Post("http://127.0.0.1:5000/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	
	// show the data sumbited
	fmt.Println("\nDatos ingresados:")
	fmt.Printf("Nombre: %s\n", name)
	fmt.Printf("Apellido: %s\n", lastname)
	fmt.Printf("RUT: %s\n", rut)
	fmt.Printf("Correo: %s\n", email)
	
	// return code
	if resp.Status == "201 Created" {
		fmt.Println("Registro exitoso!")
		return 0
	}else{
		return 1
	}
}