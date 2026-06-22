package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project-service/internal/errorz"

	sq "github.com/Masterminds/squirrel"

	"project-service/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	projectsTableName = "projects"

	projectsIDColumn                           = "id"
	projectsTitleColumn                        = "title"
	projectsDescriptionColumn                  = "description"
	projectsUsersColumn                        = "users"
	projectsAdminIDColumn                      = "admin_id"
	projectsNotificationSubscribersTGIDSColumn = "notification_subscribers_tg_ids"
	projectsCreatedAtColumn                    = "created_at"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) (string, error) {
	const op = "Postgres.ProjectRepository.Create"

	builder := sq.Insert(projectsTableName).
		Columns(projectsTitleColumn, projectsDescriptionColumn, projectsUsersColumn, projectsAdminIDColumn, projectsNotificationSubscribersTGIDSColumn).
		Values(project.Title, project.Description, project.Users, project.AdminID, project.NotificationSubscribersTGIDS).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var id string
	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	const op = "Postgres.ProjectRepository.GetByID"

	builder := sq.Select(projectsIDColumn, projectsTitleColumn, projectsDescriptionColumn, projectsUsersColumn, projectsAdminIDColumn, projectsNotificationSubscribersTGIDSColumn, projectsCreatedAtColumn).
		From(projectsTableName).
		Where(sq.Eq{projectsIDColumn: id}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var project models.RepoProject
	err = r.db.GetContext(ctx, &project, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorz.ErrProjectNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return models.FromRepoToService(&project), nil
}

func (r *ProjectRepository) GetAllByUserID(ctx context.Context, userID string) ([]*models.Project, error) {
	const op = "Postgres.ProjectRepository.GetAllByUserID"

	builder := sq.Select(projectsIDColumn, projectsTitleColumn, projectsDescriptionColumn, projectsUsersColumn, projectsAdminIDColumn, projectsNotificationSubscribersTGIDSColumn, projectsCreatedAtColumn).
		From(projectsTableName).
		Where(sq.Expr("? = ANY(users)", userID)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var repoProjects []*models.RepoProject
	err = r.db.SelectContext(ctx, &repoProjects, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	projects := make([]*models.Project, len(repoProjects))
	for i, repoProject := range repoProjects {
		projects[i] = models.FromRepoToService(repoProject)
	}

	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	const op = "Postgres.ProjectRepository.Update"

	builder := sq.Update(projectsTableName).
		Set(projectsTitleColumn, project.Title).
		Set(projectsDescriptionColumn, project.Description).
		Set(projectsUsersColumn, project.Users).
		Set(projectsAdminIDColumn, project.AdminID).
		Set(projectsNotificationSubscribersTGIDSColumn, project.NotificationSubscribersTGIDS).
		Where(sq.Eq{projectsIDColumn: project.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errorz.ErrProjectNotFound)
	}

	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	const op = "Postgres.ProjectRepository.Delete"

	builder := sq.Delete(projectsTableName).
		Where(sq.Eq{projectsIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, errorz.ErrProjectNotFound)
	}

	return nil
}
