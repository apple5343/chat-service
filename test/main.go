package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"test/client"
	"time"
)

var (
	c   *client.Client
	now int64
)

type User struct {
	id          string
	info        client.UserInfo
	accessToken string
}

func main() {
	c = client.NewClient("http://localhost:8090/api")
	now = time.Now().Unix()
	users := CreateUsers(100)
	userProjects := CreateProjects(users)
	fmt.Println(userProjects)
}

func CreateUsers(count int) []User {
	result := make([]User, count)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := range count {
		go func() {
			defer wg.Done()
			userInfo := client.UserInfo{
				Email:    fmt.Sprintf("test.%d.%d@test.com", now, i),
				Password: "qwerty1234",
				Grade:    "Test"}
			id, err := c.RegisterUser(userInfo)
			if err != nil {
				panic(err)
			}
			refreshToken, err := c.LoginUser(userInfo)
			if err != nil {
				panic(err)
			}
			accessToken, err := c.GetAccessToken(refreshToken)
			if err != nil {
				panic(err)
			}
			result[i] = User{id: id, info: userInfo, accessToken: accessToken}
		}()
	}
	wg.Wait()
	return result
}

func CreateProjects(users []User) [][]string {
	userProjects := make([][]string, len(users))
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for i := range users {
		go func() {
			defer wg.Done()
			members := make([]string, 4)
			membersIndxs := make([]int, 4)
			for j := range members {
				indx := (i + rand.IntN(len(users)-1) + 1) % len(users)
				members[j] = users[indx].id
				membersIndxs[j] = indx
			}
			project := client.Project{
				Title:       "test",
				Description: "Test",
				Users:       []string{},
			}
			fmt.Println(project)
			id, err := c.CreateProject(users[i].accessToken, project)
			if err != nil {
				panic(err)
			}
			mu.Lock()
			for _, indx := range membersIndxs {
				userProjects[indx] = append(userProjects[indx], id)
			}
			userProjects[i] = append(userProjects[i], id)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return userProjects
}
