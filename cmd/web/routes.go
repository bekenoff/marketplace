package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	// Clients
	mux.Post("/api/create-client", dynamicMiddleware.ThenFunc(app.signupClient))
	mux.Post("/api/login", dynamicMiddleware.ThenFunc(app.loginClient))

	// Products
	mux.Get("/products", dynamicMiddleware.ThenFunc(app.getProducts))           // work
	mux.Post("/api/product/add", dynamicMiddleware.ThenFunc(app.createProduct)) // work

	// Ratings
	mux.Post("/api/review/add", dynamicMiddleware.ThenFunc(app.addReview))
	mux.Get("/api/product/:id", standardMiddleware.ThenFunc(app.getProductWithRating))

	// Fav

	mux.Post("/add-favorites", dynamicMiddleware.ThenFunc(app.addFavorite))  // work
	mux.Get("/get-favorites", standardMiddleware.ThenFunc(app.getFavorites)) // work http://localhost:4000/getFavorites?id=1

	// Information

	mux.Post("/add-details", dynamicMiddleware.ThenFunc(app.addInformation))
	mux.Get("/get-details", standardMiddleware.ThenFunc(app.getInformation))

	// Image

	mux.Post("/add-image", dynamicMiddleware.ThenFunc(app.addImage))

	// Admin

	mux.Post("/api/login", dynamicMiddleware.ThenFunc(app.loginClient))

	// Order

	mux.Post("/api/add-order", dynamicMiddleware.ThenFunc(app.insertOrder))
	mux.Get("/api/get-order", standardMiddleware.ThenFunc(app.getOrderById))

	mux.Post("/api/add-order-item", dynamicMiddleware.ThenFunc(app.insertOrderItem))
	mux.Get("/api/get-order-item", standardMiddleware.ThenFunc(app.getOrdeItemrById))

	mux.Put("/api/update-status", dynamicMiddleware.ThenFunc(app.updateStatusByUserID))

	return standardMiddleware.Then(mux)
}
