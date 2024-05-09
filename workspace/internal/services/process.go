package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/T1-SD/internal/model"
)

type Files struct {
	Server_filename string `bson:"server_filename" json:"server_filename"`
	Filename        string `bson:"filename" json:"filename"`
}

type ProtectRequest struct {
	Task     string  `bson:"task" json:"task"`
	Tool     string  `bson:"tool" json:"tool"`
	Files    []Files `bson:"files" json:"files"`
	Password string  `bson:"password" json:"password"`
}

func Protect(token string, server string, task string, server_filename string, file_name string, id string) error {
	iLovePDFEndpoint := "https://" + server + "/v1/process"

	files := []Files{
		{
			Server_filename: server_filename,
			Filename:        file_name,
		},
	}

	password, err := ObtainPassword(id)
	if err != nil {
		return err
	}

	body := &ProtectRequest{
		Task:     task,
		Tool:     "protect",
		Files:    files,
		Password: password,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error al convertir el cuerpo a json: %v", err)
	}

	req, err := http.NewRequest("POST", iLovePDFEndpoint, strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("error creando la solicitud: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta no fue exitosa: %d", resp.StatusCode)
	}

	return nil
}

func ObtainPassword(id string) (string, error) {
	resp, err := http.Get(os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/clients/" + id)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error al obtener el cliente")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var client model.Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		return "", err
	}

	password := strings.Split(client.Rut, "-")[0]

	return password, nil
}
