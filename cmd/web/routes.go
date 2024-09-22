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
	mux.Get("/api/get-product", dynamicMiddleware.ThenFunc(app.getProductByID))

	// Authenticated routes (JWT required)
	mux.Post("/api/create-client", authMiddleware.ThenFunc(app.signupClient))
	mux.Post("/api/add-order", authMiddleware.ThenFunc(app.insertOrder))
	mux.Get("/api/get-order", authMiddleware.ThenFunc(app.getOrderById))
	mux.Get("/api/cart", authMiddleware.ThenFunc(app.getCartItems))
	mux.Post("/api/cart/add", authMiddleware.ThenFunc(app.addCartItem))
	mux.Del("/api/cart/remove", authMiddleware.ThenFunc(app.deleteCartItem))

	// Apply to dynamic routes as needed
	mux.Post("/api/product/add", authMiddleware.ThenFunc(app.createProduct))
	mux.Post("/add-favorites", authMiddleware.ThenFunc(app.addFavorite))

	// Admin routes can also use this
	mux.Put("/api/update-status", authMiddleware.ThenFunc(app.updateStatusByUserID))

	// Return the final handler with all standard middleware
	return standardMiddleware.Then(mux)
}
