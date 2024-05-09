package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UploadResponse struct {
	Server_filename string `bson:"server_filename" json:"server_filename"`
}

func Upload(token string, server string, body *bytes.Buffer) (string, error) {
	iLovePDFEndpoint := "https://" + server + "/v1/start/upload"

	req, err := http.NewRequest("POST", iLovePDFEndpoint, body)

	if err != nil {
		return "", fmt.Errorf("error creando la solicitud: %v", err)
	}

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
