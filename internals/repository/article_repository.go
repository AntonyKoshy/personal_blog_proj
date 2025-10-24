package repository

import (
	"context"
	"fmt"
	"personal-blog/internals/config"
	"personal-blog/internals/models"

	"github.com/jackc/pgx/v5"
)

type ArticleRepository struct {
	Conn *pgx.Conn
}

// func to connect to DB
func ConnectDB(cfg *config.Config) (*pgx.Conn, error) {

	dsn := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return pgx.Connect(context.Background(), dsn)

}

// Fetch all Articles
func (r *ArticleRepository) FetchArticles(ctx context.Context) ([]models.Article, error) {

	query := `SELECT id, content, title , created_at, updated_at, category, tags FROM articles`
	rows, err := r.Conn.Query(context.Background(), query)
	if err != nil {

		return nil, err

	}

	defer rows.Close()
	var articles []models.Article
	for rows.Next() {

		var a models.Article
		if err := rows.Scan(&a.ID, &a.Content, &a.Title, &a.CreatedAt, &a.UpdatedAt, &a.Category, &a.Tags); err != nil {
			return nil, pgx.ErrNoRows

		}

		articles = append(articles, a)
	}

	return articles, nil

}

// Insert Article
func (r *ArticleRepository) InsertArticles(ctx context.Context, a *models.Article) error {

	query := `INSERT INTO articles (title, content, category, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.Conn.QueryRow(ctx, query, a.Title, a.Content, a.Category, a.Tags).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

}

// Fetch Articles by ID
func (r *ArticleRepository) FetchArticleByID(ctx context.Context, id int) (models.Article, error) {

	query := `SELECT id, title, content, created_at, updated_at, category, tags FROM articles WHERE id = $1`
	var a models.Article
	err := r.Conn.QueryRow(ctx, query, id).Scan(&a.ID, &a.Title, &a.Content, &a.CreatedAt, &a.UpdatedAt, &a.Category, &a.Tags)
	return a, err

}
