package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ClientMainMenu() {
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
		fmt.Println(menu_clients)

		option := ReadOption()
		switch option {
		case 1:
			getAllClients()
		case 2:
			fmt.Print("\nIngrese el id a buscar: ")
			var id string
			fmt.Scan(&id)
			getClientByID(id)
		case 3:
			fmt.Println("\nObtener cliente por RUT no implementado.")
		case 4:
			createClient()
		case 5:
			fmt.Print("\nIngrese el id del cliente a modificar: ")
			var id string
			fmt.Scan(&id)
			updateClient(id)
		case 6:
			fmt.Print("\nIngrese el id del cliente a eliminar: ")
			var id string
			fmt.Scan(&id)
			deleteClient(id)
		case 7:
			fmt.Println("\nVolviendo al menú principal...")
			return
		default:
			fmt.Println("\nOpción no válida. Intenta de nuevo.")
		}
	}
}

func getAllClients() {
	resp, err := http.Get(os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/clients")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Error al obtener los clientes")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n---")
	fmt.Println(string(body))
	fmt.Println("---")
}

func getClientByID(id string) {
	resp, err := http.Get(os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/clients/" + id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Error al obtener el cliente")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n---")
	fmt.Println(string(body))
	fmt.Println("---")
}

func createClient() {
	fmt.Print("\nIngrese nombre del cliente: ")
	var name string
	fmt.Scan(&name)

	fmt.Print("\nIngrese apellido: ")
	var lastname string
	fmt.Scan(&lastname)

	fmt.Print("\nIngrese RUT del cliente: ")
	var rut string
	fmt.Scan(&rut)

	fmt.Print("\nIngrese el correo del cliente: ")
	var mail string
	fmt.Scan(&mail)

	body := map[string]interface{}{
		"Name":      name,
		"Last_name": lastname,
		"Rut":       rut,
		"Email":     mail,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/api/clients", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("\n¡Cliente " + name + " creado con éxito!")
}

func updateClient(id string) {
	fmt.Print("\nIngrese el nuevo nombre del cliente: ")
	var name string
	fmt.Scan(&name)

	fmt.Print("\nIngrese el nuevo apellido: ")
	var lastname string
	fmt.Scan(&lastname)

	fmt.Print("\nIngrese el nuevo RUT del cliente: ")
	var rut string
	fmt.Scan(&rut)

	fmt.Print("\nIngrese el nuevo correo del cliente: ")
	var mail string
	fmt.Scan(&mail)

	body := map[string]interface{}{
		"Name":      name,
		"Last_name": lastname,
		"Rut":       rut,
		"Email":     mail,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, os.Getenv("HOST")+":"+os.Getenv("PORT")+"/api/clients/"+id, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("\n¡Cliente " + name + " modificado con éxito!")
}

func deleteClient(id string) {
	req, err := http.NewRequest(http.MethodDelete, os.Getenv("HOST")+":"+os.Getenv("PORT")+"/api/clients/"+id, nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("\n¡Cliente eliminado con éxito!")
}
