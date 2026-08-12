package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sort"
	"sync"
	"test/client"
	"test/rps"
	"test/storage"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
)

var (
	c          *client.Client
	rpsCounter *rps.RPSCounter
	now        int64
)

type User struct {
	id          string
	info        client.UserInfo
	accessToken string
}

func main() {
	c = client.NewClient("http://localhost:8090/api", "ws://localhost:8090/api")
	rpsCounter = rps.NewRPSCounter()
	now = time.Now().Unix()
	fmt.Println("Create users...")
	users := CreateUsers(1000)
	userProjects := CreateProjects(users)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rpsCounter.Run(ctx)
	fmt.Println("Send messages...")
	start := time.Now()
	messages := SendMessages(ctx, users, userProjects, 60)
	cancel()
	fmt.Println(rpsCounter.History())
	Analyze(messages, start)
}

func CreateUsers(count int) []User {
	result := make([]User, count)
	wg := sync.WaitGroup{}
	wg.Add(count)
	sem := semaphore.NewWeighted(100)
	for i := range count {
		go func() {
			defer wg.Done()
			sem.Acquire(context.Background(), 1)
			defer sem.Release(1)
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
	sem := semaphore.NewWeighted(100)
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
				Users:       members,
			}
			sem.Acquire(context.Background(), 1)
			defer sem.Release(1)
			id, err := c.CreateProject(users[i].accessToken, project)
			if err != nil {
				panic(err)
			}
			if id == "" {
				panic("empty id")
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

func SendMessages(ctx context.Context, users []User, projects [][]string, count int) [][]storage.MessagePath {
	messages := make([][]storage.MessagePath, len(users))

	wg := sync.WaitGroup{}
	wg.Add(len(users))
	for i, user := range users {
		go func() {
			defer wg.Done()
			conn, err := c.ConnectWS(ctx, user.accessToken)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			time.Sleep(time.Duration(rand.Int64N(2000)) * time.Millisecond)
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			storage := storage.NewStorage()

			go func() {
				type body struct {
					Content string `json:"content"`
				}
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						return
					}
					now := time.Now()
					var b body
					if err := json.Unmarshal(message, &b); err != nil {
						continue
					}
					storage.Receive(b.Content, now)
				}
			}()

			for range count {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					id := uuid.NewString()
					storage.Send(id, time.Now())
					projectID := projects[i][rand.IntN(len(projects[i]))]
					rpsCounter.Inc()
					if err := c.SendMessage(conn, &client.Message{RoomID: projectID, Content: id}); err != nil {
						panic(err)
					}
				}
			}

			time.Sleep(time.Second)
			conn.Close()
			messages[i] = storage.Messages()
		}()
	}
	wg.Wait()
	return messages
}

func Analyze(messages [][]storage.MessagePath, start time.Time) {
	sended := 0
	received := 0
	from := start.Add(time.Second * 5)
	var totalDur time.Duration
	for _, msgs := range messages {
		for _, msg := range msgs {
			if msg.SendedAt.Before(from) {
				continue
			}
			sended++
			if msg.Received {
				received++
				totalDur += msg.ReceivedAt.Sub(msg.SendedAt)
			}
		}
	}
	fmt.Println("Sended: ", sended)
	fmt.Println("Received: ", received)
	fmt.Println("Avg latency: ", totalDur/time.Duration(received))
	Percentiles(messages, from)
}

func Percentiles(messages [][]storage.MessagePath, from time.Time) {
	var durations []time.Duration
	for _, msgs := range messages {
		for _, msg := range msgs {
			if msg.SendedAt.Before(from) {
				continue
			}
			if msg.Received {
				durations = append(durations, msg.ReceivedAt.Sub(msg.SendedAt))
			}
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p := func(pct float64) time.Duration {
		idx := int(float64(len(durations)-1) * pct)
		return durations[idx]
	}
	fmt.Println("p50:", p(0.50), "p95:", p(0.95), "p99:", p(0.99))
}
