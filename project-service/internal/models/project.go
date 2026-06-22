package models

import (
	"time"

	"github.com/lib/pq"
)

type Project struct {
	ID                           string    `json:"id"`
	Title                        string    `json:"title"`
	Description                  string    `json:"description"`
	Users                        []string  `json:"users"`
	AdminID                      string    `json:"admin_id"`
	NotificationSubscribersTGIDS []int64   `json:"notification_subscribers_tg_ids"`
	CreatedAt                    time.Time `json:"created_at" `
}

type RepoProject struct {
	ID                           string         `db:"id"`
	Title                        string         `db:"title"`
	Description                  string         `db:"description"`
	Users                        pq.StringArray `db:"users"`
	AdminID                      string         `db:"admin_id"`
	NotificationSubscribersTGIDS pq.Int64Array  `db:"notification_subscribers_tg_ids"`
	CreatedAt                    time.Time      `db:"created_at"`
}

func FromRepoToService(repoProject *RepoProject) *Project {
	return &Project{
		ID:                           repoProject.ID,
		Title:                        repoProject.Title,
		Description:                  repoProject.Description,
		Users:                        repoProject.Users,
		AdminID:                      repoProject.AdminID,
		NotificationSubscribersTGIDS: repoProject.NotificationSubscribersTGIDS,
		CreatedAt:                    repoProject.CreatedAt,
	}
}
