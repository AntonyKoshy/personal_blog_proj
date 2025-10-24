package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"personal-blog/internals/models"
	"personal-blog/internals/repository"
	"strconv"
)

func GetArticlesHandler(repo *repository.ArticleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		articles, err := repo.FetchArticles(context.Background())
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

		err := repo.InsertArticles(context.Background(), &a)
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

		article, err := repo.FetchArticleByID(context.Background(), id)
		if err != nil {
			http.Error(w, "Failed to fetch article", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)

	}
}
