package main

import (
	"encoding/json"
	"errors"
	"marketplace/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) insertOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order

	// Decode the incoming request body into the Order struct
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Insert the order and update inventory
	err = app.order.Insert(newOrder.User_id, newOrder.Status, newOrder.Address, newOrder.Price, newOrder.Product_id, newOrder.Quantity)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getOrderById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	order, err := app.order.GetOrderById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(order)
}

func (app *application) updateStatusByUserID(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID int    `json:"user_id"`
		Status string `json:"status"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.order.UpdateStatusByUserID(req.UserID, req.Status)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) insertOrderItem(w http.ResponseWriter, r *http.Request) {
	var newOrder models.OrderItem

	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.order.InsertOrderItem(newOrder.OrderID, newOrder.ProductID, newOrder.Price, newOrder.Qty)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getOrdeItemrById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	order, err := app.order.GetOrderItemById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(order)
}
