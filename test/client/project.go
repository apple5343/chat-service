package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
}

func (c *Client) CreateProject(accessToken string, project Project) (string, error) {
	b, err := json.Marshal(project)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/projects", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}
	return respBody.ID, nil
}
