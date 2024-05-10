package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/T1-SD/internal/model"
	"github.com/T1-SD/internal/services"
)

func ProtectFile(c *gin.Context) {

	token, err := GetTokenString()
	if err != nil {
		fmt.Println("Error obteniendo el token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo el token"})
		return
	}

	server, task, err := services.Start(token)
	if err != nil {
		fmt.Println("Error en el servicio Start:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio Start"})
		return
	}

	fmt.Println("server:", server)
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println("task:", task)

	clientID := c.PostForm("id")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error obteniendo el archivo del form-data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no encontrado"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo"})
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("task", task)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		fmt.Println("Error creando el form-data para el archivo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando el form-data"})
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error al copiar el archivo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al copiar el archivo"})
		return
	}

	writer.Close()
	reqContentType := writer.FormDataContentType()

	fmt.Println("token:", token)
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA/nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println("id:", clientID)
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA/nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	server_filename, err := services.Upload(token, server, body, reqContentType)
	if err != nil {
		fmt.Println("Error en el servicio Upload:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio Upload"})
		return
	}

	fmt.Println("\nCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC\nCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC\nCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	fmt.Println("fileHeader.Filename:", fileHeader.Filename)
	fmt.Println("server_filename:", server_filename)

	err = services.Protect(token, server, task, server_filename, fileHeader.Filename, clientID)
	if err != nil {
		fmt.Println("Error en el servicio Protect:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio Protect"})
		return
	}

	downloaded_body, err := services.Download(token, server, task)
	if err != nil {
		fmt.Println("Error en el servicio Download:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio Download"})
		return
	}

	fmt.Print("\n\nAAAAAAA\ndownloaded_body: ")
	fmt.Println(downloaded_body)
	fmt.Print("AAAAAAAA\n\n")

	filename := fileHeader.Filename + "_protected.pdf"
	outputDir := "/pdfs"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creando el directorio:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando el directorio"})
			return
		}
	}

	filePath := filepath.Join(outputDir, filename)
	err = os.WriteFile(filePath, downloaded_body, 0644)
	if err != nil {
		fmt.Println("Error guardando el archivo descargado:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error guardando el archivo descargado"})
		return
	}

	err = services.Delete(token, server, task, server_filename)
	if err != nil {
		fmt.Println("Error en el servicio Delete:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio Delete"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "pdf protected successfully"})
}

func GetTokenString() (string, error) {
	resp, err := http.Get(os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/token")
	if err != nil {
		fmt.Println("Error al hacer la solicitud GET:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error al obtener el token")
		return "", fmt.Errorf("error al obtener el token: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return "", err
	}

	var tokenResp model.CreateTokenRequest
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("error al decodificar el JSON: %v", err)
	}

	return tokenResp.Token, nil
}
