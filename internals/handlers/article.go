package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"personal-blog/internals/models"
	"personal-blog/internals/repository"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func GetArticlesHandler(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		searchTerm := r.URL.Query().Get("post")

		var articles []models.Article
		var err error

		if searchTerm == "" {
			articles, err = repo.FetchArticles(r.Context())
		} else {
			articles, err = repo.FetchArticlesByTerm(r.Context(), searchTerm)

		}

		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)

	}
}

func CreateArticleHandler(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var a models.Article
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := repo.InsertArticles(r.Context(), &a)
		if err != nil {
			http.Error(w, "Failed to insert article", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(a)

	}
}

// Handler to parse out id from URL & invoke repo func
func GetArticlesByIDHandler(repo *repository.ArticleRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		article, err := repo.FetchArticleByID(r.Context(), id)
		if err != nil {

			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Article not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "Failed to fetch article", http.StatusInternalServerError)
				return
			}

		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)

	}
}

func UpdateArticleHandler(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return

		}

		var a models.Article
		if err = json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return

		}

		err = repo.UpdateArticle(r.Context(), id, &a)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to update article", http.StatusInternalServerError)
			return

		}
		a.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)

	}

}

func DeleteArticleHandler(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return

		}

		err = repo.DeleteArticle(r.Context(), id)

		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to delete article", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}

}

func GetArticlesByTerm(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			return
		}
		term := r.URL.Query().Get("term")
		articles, err := repo.FetchArticlesByTerm(r.Context(), term)
		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
}
