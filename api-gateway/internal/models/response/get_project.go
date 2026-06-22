package response

import "time"

type GetProject struct {
	ID                           string    `json:"id"`
	Title                        string    `json:"title"`
	Description                  string    `json:"description"`
	Users                        []string  `json:"users"`
	AdminID                      string    `json:"admin_id"`
	NotificationSubscribersTGIds []int64   `json:"notification_subscribers_tg_ids"`
	CreatedAt                    time.Time `json:"created_at"`
}
