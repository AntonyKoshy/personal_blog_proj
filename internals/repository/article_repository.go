package repository

import (
	"context"
	"errors"
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

	query := `SELECT id, title, content, created_at, updated_at, category, tags FROM articles ORDER BY created_at ASC`
	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Article{}, nil

		} else {
			return nil, err

		}
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Article])
	if err != nil {
		return nil, err
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

func (r *ArticleRepository) UpdateArticle(ctx context.Context, id int, a *models.Article) error {

	query := `UPDATE articles SET title = $1, content = $2, category = $3, tags = $4 WHERE id = $5`
	cmdTag, err := r.Conn.Exec(ctx, query, a.Title, a.Content, a.Category, a.Tags, id)

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err

}

func (r *ArticleRepository) DeleteArticle(ctx context.Context, id int) error {

	query := `DELETE FROM articles WHERE id = $1`

	cmdTag, err := r.Conn.Exec(ctx, query, id)
	fmt.Println(cmdTag)

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err

}

func (r *ArticleRepository) FetchArticlesByTerm(ctx context.Context, post string) ([]models.Article, error) {

	query := `
		SELECT id, title, content, created_at, updated_at, category, tags 
		FROM articles 
		WHERE title ILIKE $1 
		   OR content ILIKE $1 
		   OR category ILIKE $1 
		   OR array_to_string(tags, ' ') ILIKE $1`

	searchTerm := "%" + post + "%"
	rows, err := r.Conn.Query(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Article])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Article{}, nil // Return empty slice if no articles found
		}
		return nil, err
	}
	return articles, nil

}
