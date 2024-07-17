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

	// Sights
	mux.Post("/api/sight/add", dynamicMiddleware.ThenFunc(app.createSight))       // work
	mux.Get("/api/sight/get/:id", standardMiddleware.ThenFunc(app.getSight))      //http://localhost:4000/api/sight/get/2 work
	mux.Get("/api/sights/get-all", standardMiddleware.ThenFunc(app.getAllSights)) // work

	// Events
	mux.Post("/api/event/add", dynamicMiddleware.ThenFunc(app.createEvent))
	mux.Get("/api/event/get/:id", standardMiddleware.ThenFunc(app.getEvent))     //http://localhost:4000/api/event/get/2 work
	mux.Get("/api/event/get-all", standardMiddleware.ThenFunc(app.GetAllEvents)) // work

	// Rec
	mux.Post("/api/rec/add", dynamicMiddleware.ThenFunc(app.createRec))     // workk
	mux.Get("/api/rec/get/:id", standardMiddleware.ThenFunc(app.getRec))    //http://localhost:4000/api/rec/get/2
	mux.Get("/api/rec/get-all", standardMiddleware.ThenFunc(app.GetAllRec)) // work

	// Fav

	mux.Post("/addFavorites", dynamicMiddleware.ThenFunc(app.addFavorite))
	mux.Get("/getFavorites", standardMiddleware.ThenFunc(app.getFavorites))

	return standardMiddleware.Then(mux)
}
