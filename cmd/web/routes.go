package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Standard middleware (security, logging, panic recovery)
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Middleware for session management
	dynamicMiddleware := alice.New(app.session.Enable)

	// Middleware for JWT authentication (e.g., for protected routes)
	authMiddleware := alice.New(app.jwtAuthMiddleware)

	mux := pat.New()

	// Public routes (no authentication required)
	mux.Post("/api/login", dynamicMiddleware.ThenFunc(app.loginClient))

	mux.Post("/api/create-client-law", dynamicMiddleware.ThenFunc(app.signupClientLaw))
	mux.Post("/api/login-admin", dynamicMiddleware.ThenFunc(app.loginAdmin))
	mux.Get("/api/get-product", dynamicMiddleware.ThenFunc(app.getProductByID))

	// Authenticated routes (JWT required)
	mux.Post("/api/create-client", dynamicMiddleware.ThenFunc(app.signupClient))
	mux.Get("/api/user-by-id", dynamicMiddleware.ThenFunc(app.getUserById))
	mux.Post("/api/information", dynamicMiddleware.ThenFunc(app.addInformation))
	mux.Get("/api/information", dynamicMiddleware.ThenFunc(app.getInformation))
	mux.Put("/api/update-password", dynamicMiddleware.ThenFunc(app.updatePassword))
	mux.Post("/api/add-order", authMiddleware.ThenFunc(app.insertOrder))
	mux.Get("/api/get-order", authMiddleware.ThenFunc(app.getOrderById))
	mux.Get("/api/cart", authMiddleware.ThenFunc(app.getCartItems))
	mux.Post("/api/cart/add", authMiddleware.ThenFunc(app.addCartItem))
	mux.Del("/api/cart/remove", authMiddleware.ThenFunc(app.deleteCartItem))

	// Apply to dynamic routes as needed
	mux.Post("/api/product/add", authMiddleware.ThenFunc(app.createProduct))
	mux.Post("/add-favorites", authMiddleware.ThenFunc(app.addFavorite))
	mux.Get("/get-favorites", authMiddleware.ThenFunc(app.getFavorites))
	mux.Post("/refresh-token", authMiddleware.ThenFunc(app.refreshToken))

	// Admin routes can also use this
	mux.Put("/api/update-status", authMiddleware.ThenFunc(app.updateStatusByUserID))

	// Discount

	mux.Post("/api/discount", authMiddleware.ThenFunc(app.addDiscount))

	// Image

	mux.Post("/api/image", authMiddleware.ThenFunc(app.addImage))

	// OrderItem

	mux.Post("/api/order-item", authMiddleware.ThenFunc(app.insertOrderItem))
	mux.Get("/api/order-item", authMiddleware.ThenFunc(app.getOrdeItemrById))

	// Products

	mux.Get("/api/products", authMiddleware.ThenFunc(app.getProducts))

	mux.Get("/api/products-with-rating", authMiddleware.ThenFunc(app.getProductsWithRating))

	// Review
	mux.Post("/api/add-review", authMiddleware.ThenFunc(app.addReview))
	mux.Get("/api/product-rating", authMiddleware.ThenFunc(app.getProductWithRating))

	mux.Post("/api/product-inventory", authMiddleware.ThenFunc(app.createProductInventory))
	mux.Get("/api/product-by-category-id", authMiddleware.ThenFunc(app.getProductByCategoryID))

	// Return the final handler with all standard middleware
	return standardMiddleware.Then(mux)
}
