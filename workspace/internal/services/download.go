package services

import (
	"fmt"
	"io"
	"net/http"
)

func Download(token string, server string, task string) ([]byte, error) {
	iLovePDFEndpoint := "https://" + server + "/v1/download/" + task

	req, err := http.NewRequest("GET", iLovePDFEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta no fue exitosa: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
