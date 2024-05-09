package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StartResponse struct {
	Server string `bson:"server" json:"server"`
	Task   string `bson:"task" json:"task"`
}

func Start(token string) (string, string, error) {
	iLovePDFEndpoint := "https://api.ilovepdf.com/v1/start/protect"

	req, err := http.NewRequest("GET", iLovePDFEndpoint, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("respuesta no fue exitosa: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var startResp StartResponse
	err = json.Unmarshal(body, &startResp)
	if err != nil {
		return "", "", err
	}

	return startResp.Server, startResp.Task, nil
}
