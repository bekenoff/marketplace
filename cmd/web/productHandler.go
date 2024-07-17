package main

import (
	"encoding/json"
	"errors"
	"marketplace/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func getProducts(c *gin.Context) {
	var products []models.Product
	db.Find(&products)

	var productsWithRating []models.ProductWithRating
	for _, product := range products {
		var reviews []models.Review
		db.Where("product_id = ?", product.ID).Find(&reviews)

		var totalRating int
		for _, review := range reviews {
			totalRating += review.Rating
		}

		averageRating := 0.0
		if len(reviews) > 0 {
			averageRating = float64(totalRating) / float64(len(reviews))
		}

		productsWithRating = append(productsWithRating, models.ProductWithRating{
			Product:       product,
			AverageRating: averageRating,
		})
	}

	c.JSON(http.StatusOK, productsWithRating)
}

func createReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&review)
	c.JSON(http.StatusCreated, review)
}

func (app *application) addReview(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductID int    `json:"product_id"`
		Rating    int    `json:"rating"`
		Review    string `json:"review"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Валидация данных
	if input.Rating < 1 || input.Rating > 5 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.product.InsertRating(input.ProductID, input.Rating, input.Review)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getProductWithRating(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	product, err := app.product.GetProductByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	reviews, err := app.product.GetReviewsByProductID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}

	averageRating := 0.0
	if len(reviews) > 0 {
		averageRating = float64(totalRating) / float64(len(reviews))
	}

	productWithRating := struct {
		Product       models.Product `json:"product"`
		AverageRating float64        `json:"average_rating"`
	}{
		Product:       *product,
		AverageRating: averageRating,
	}

	err = json.NewEncoder(w).Encode(productWithRating)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err := app.product.InsertProduct(product.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
