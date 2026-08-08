package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Grade    string `json:"grade"`
}

func (c *Client) RegisterUser(userInfo UserInfo) (string, error) {
	b, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/register", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
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

func (c *Client) LoginUser(userInfo UserInfo) (string, error) {
	b, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/login", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}
	return respBody.RefreshToken, nil
}

func (c *Client) GetAccessToken(refreshToken string) (string, error) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	body.RefreshToken = refreshToken
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/get_access_token", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var respBody struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}
	return respBody.AccessToken, nil
}
