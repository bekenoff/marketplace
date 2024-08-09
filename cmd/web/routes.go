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
	mux.Post("/api/create-client-law", dynamicMiddleware.ThenFunc(app.signupClientLaw))
	mux.Get("/api/get-client", dynamicMiddleware.ThenFunc(app.getUserById))
	mux.Post("/api/login", dynamicMiddleware.ThenFunc(app.loginClient))
	mux.Put("/client-password-recovery/:id", dynamicMiddleware.ThenFunc(app.Recoverybysms))

	// Products
	mux.Get("/all-products", dynamicMiddleware.ThenFunc(app.getProducts)) // work
	mux.Get("/api/get-product", dynamicMiddleware.ThenFunc(app.getProductByID))
	mux.Get("/api/get-product-inventory", dynamicMiddleware.ThenFunc(app.getProductByCategoryID))
	mux.Post("/api/product/add", dynamicMiddleware.ThenFunc(app.createProduct)) // work
	mux.Post("/api/product-inventory/add", dynamicMiddleware.ThenFunc(app.createProductInventory))

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

	// Cart
	mux.Get("/api/cart", dynamicMiddleware.ThenFunc(app.getCartItems))
	mux.Post("/api/cart/add", dynamicMiddleware.ThenFunc(app.addCartItem))
	mux.Del("/api/cart/remove", dynamicMiddleware.ThenFunc(app.deleteCartItem))

	return standardMiddleware.Then(mux)
}
