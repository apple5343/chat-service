package repository

import (
	"context"
	"user-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// почему интерфейс по месту объявления а не использования?

// согласовали хранить их по месту объявления. Предложил идею я (Юра, в тг @uikola). До недавних пор придерживался мнения
// Николая Тузова, что храним по месту использования, но сейчас т.к. узнал о DI контейнере, столкнулся довольно с неприятной вещью.
// Если мы храним по месту использования, то в пакет с DI контейнером приходится вывалить абсолютно все интерфейсы, а если мы храним
// по месту объявления, то мы просто импортируем их из пакета с интерфейсами для того или иного слоя.
