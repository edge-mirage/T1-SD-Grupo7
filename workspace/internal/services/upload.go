package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UploadResponse struct {
	Server_filename string `bson:"server_filename" json:"server_filename"`
}

func Upload(token string, server string, task string, file multipart.File, file_name string) (string, error) {
	iLovePDFEndpoint := "https://" + server + "/v1/start/upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("task", task)
	if err != nil {
		return "", fmt.Errorf("error agregando task_id al cuerpo: %v", err)
	}

	part, err := writer.CreateFormFile("file", file_name)
	if err != nil {
		return "", fmt.Errorf("error creando parte para el archivo: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copiando el archivo: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error cerrando el writer: %v", err)
	}

	req, err := http.NewRequest("POST", iLovePDFEndpoint, body)
	if err != nil {
		return "", fmt.Errorf("error creando la solicitud: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error enviando la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("respuesta no fue exitosa: %d", resp.StatusCode)
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var uploadResp UploadResponse
	err = json.Unmarshal(resp_body, &uploadResp)
	if err != nil {
		return "", err
	}

	return uploadResp.Server_filename, nil
}
