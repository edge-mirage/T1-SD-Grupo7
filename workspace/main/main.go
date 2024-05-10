package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

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

	token, err := ObtainToken()
	if err != nil {
		fmt.Println("Error al obtener el token:", err)
		return
	}
	PostToken(token)

	for {
		option := MainMenu(menu)
		switch option {
		case 1:
			login(menu_loged)
		case 2:
			register()
		case 3:
			fmt.Println("\n¡Vuelve pronto!")
			DeleteToken()
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
	var option int
	fmt.Scan(&option)
	return option
}

func ObtainToken() (string, error) {
	requestBody := map[string]string{"public_key": os.Getenv("PUBLIC_KEY")}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://api.ilovepdf.com/v1/auth", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]string
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	token, ok := response["token"]
	if !ok {
		return "", err
	}

	return token, nil
}

func PostToken(token string) {

	requestBody, err := json.Marshal(map[string]interface{}{
		"Token": token,
	})
	if err != nil {
		fmt.Println("Error al crear el cuerpo JSON:", err)
	}

	resp, err := http.Post(os.Getenv("HOST")+":"+os.Getenv("PORT")+"/api/token", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()
}

func DeleteToken() {
	req, err := http.NewRequest("DELETE", os.Getenv("HOST")+":"+os.Getenv("PORT")+"/api/token", nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud DELETE:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error al hacer la solicitud DELETE:", err)
		return
	}
	defer resp.Body.Close()
}

func protectFile() {
	fmt.Println("\nEscriba el ID del cliente objetivo: ")
	var clientID string
	fmt.Scan(&clientID)

	fmt.Println("\nEscriba la ruta donde se encuentra el archivo (incluya el nombre): ")
	var filePath string
	fmt.Scan(&filePath)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	err := writer.WriteField("id", clientID)
	if err != nil {
		fmt.Println("Error al escribir el campo 'id':", err)
		return
	}

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error al crear parte para el archivo:", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error al copiar el archivo:", err)
		return
	}

	writer.Close()

	url := os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/protect"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		fmt.Println("Error al crear la solicitud POST:", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al enviar la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		fmt.Printf("Error: %d\n", resp.StatusCode)
	} else {
		fmt.Println("¡Protección exitosa!")
	}
}
