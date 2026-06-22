package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	BASE_URL = "http://localhost:8090/api"
)

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Grade    string `json:"grade"`
}

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
}

func registerUser(userInfo UserInfo) (string, error) {
	b, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/register", bytes.NewBuffer(b))
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

func loginUser(userInfo UserInfo) (string, error) {
	b, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/login", bytes.NewBuffer(b))
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

func getAccessToken(refreshToken string) (string, error) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	body.RefreshToken = refreshToken
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/get_access_token", bytes.NewBuffer(b))
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

func createProject(accessToken string, project Project) (string, error) {
	b, err := json.Marshal(project)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/projects", bytes.NewBuffer(b))
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

func main() {
	users := make([]string, 0, 5)
	usersInfo := make([]UserInfo, 0, 5)
	accessTokens := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		u := UserInfo{
			Email:    fmt.Sprintf("user%d@example.com", i+1),
			Password: fmt.Sprintf("user%d", i+1),
			Grade:    "Student",
		}
		id, err := registerUser(u)
		if err != nil {
			panic(err)
		}
		users = append(users, id)
		usersInfo = append(usersInfo, u)

		refreshToken, err := loginUser(u)
		if err != nil {
			panic(err)
		}
		accessToken, err := getAccessToken(refreshToken)
		if err != nil {
			panic(err)
		}
		accessTokens = append(accessTokens, accessToken)
	}
	user2 := users[1]
	user3 := users[2]
	user4 := users[3]
	user5 := users[4]

	accessToken := accessTokens[0]

	project1 := Project{
		Title:       "Project 1",
		Description: "Description 1",
		Users:       []string{user2, user3},
	}

	project2 := Project{
		Title:       "Project 2",
		Description: "Description 2",
		Users:       []string{user4, user5},
	}

	project1ID, err := createProject(accessToken, project1)
	if err != nil {
		panic(err)
	}

	project2ID, err := createProject(accessToken, project2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Users (%d):\n", len(users))
	for i := range users {
		fmt.Printf("\t ID: %s\n", users[i])
		fmt.Printf("\t Email: %s\n", usersInfo[i].Email)
		fmt.Printf("\t Password: %s\n", usersInfo[i].Password)
		fmt.Printf("\t Grade: %s\n", usersInfo[i].Grade)
		fmt.Printf("\t Access Token: %s\n\n", accessTokens[i])
	}

	fmt.Printf("\nProjects (%d):\n", 2)

	fmt.Printf("\t ID: %s\n", project1ID)
	fmt.Printf("\t Admin: user1\n")
	fmt.Printf("\t Title: %s\n", project1.Title)
	fmt.Printf("\t Description: %s\n", project1.Description)
	fmt.Printf("\t Users: %v\n", project1.Users)
	fmt.Printf("\t Team: user1, user2, user3\n\n")

	fmt.Printf("\t ID: %s\n", project2ID)
	fmt.Printf("\t Admin: user1\n")
	fmt.Printf("\t Title: %s\n", project2.Title)
	fmt.Printf("\t Description: %s\n", project2.Description)
	fmt.Printf("\t Users: %v\n", project2.Users)
	fmt.Printf("\t Team: user1, user4, user5\n\n")

}
