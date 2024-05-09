package services

import (
	"fmt"
	"net/http"
)

func Download(token string, server string, task string) error {
	iLovePDFEndpoint := "https://" + server + "/v1/download/" + task

	req, err := http.NewRequest("GET", iLovePDFEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta no fue exitosa: %d", resp.StatusCode)
	}
	return nil
}
