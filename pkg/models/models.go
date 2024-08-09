package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Product struct {
	Id           uint    `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category_id  int     `json:"category_id"`
	Inventory_id int     `json:"inventory_id"`
	Price        float64 `json:"price"`
	Discount_id  int     `json:"discount_id"`
	Created_at   string  `json:"created_at"`
	Modified_at  string  `json:"modified_at"`
}

type Cart struct {
	Id         int `json:"id"`
	Client_id  int `json:"client_id"`
	Product_id int `json:"product_id"`
	Quantity   int `json:"quantity"`
}

type Favorites struct {
	Product_id int `json:"product_id"`
	Client_id  int `json:"client_id"`
}

type ProductWithRating struct {
	Product
	AverageRating float64 `json:"average_rating"`
}

type Review struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	Rating    int    `json:"rating"`
	Review    string `json:"review"`
}

type Client struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	First_name  string `json:"firstname"`
	Last_name   string `json:"lastname"`
	Telephone   int    `json:"telephone"`
	Created_at  string `json:"created_at"`
	Modified_at string `json:"modified_at"`
}

type Client_Events struct {
	Id       int `json:"id"`
	ClientId int `json:"client_id"`
	EventId  int `json:"event_id"`
}

type Information struct {
	Id           int `json:"id"`
	Product_id   int `json"product_id"`
	Articul      int `json"articul"`
	Brand        int `json"brand"`
	Series       int `json"series"`
	Country      int `json"country"`
	Color        int `json"color"`
	Quantity     int `json"quantity"`
	Size         int `json"size"`
	Packing_size int `json"packing_size"`
}

type Image struct {
	Id         int    `json:"id"`
	Product_id int    `json:"product_id"`
	Image_url  string `json:"image_url"`
}

type Order struct {
	Id      int    `json:"id"`
	User_id int    `json:"user_id"`
	Status  string `json:"status"`
	Address string `json:"address"`
	Price   int    `json:"price"`
}

type OrderItem struct {
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Price     int `json:"price"`
	Qty       int `json:"qty"`
}

type ClientLaw struct {
	CompanyName string `json:"company_name"`
	ContactName string `json:"contact_name"`
	Password    string `json:"password"`
	LawAddress  string `json:"law_address"`
	Email       string `json:"email"`
	Phone       int    `json:"phone"`
	Bin         int    `json:"bin"`
	Bik         int    `json:"bik"`
	Iik         int    `json:"iik"`
	Bank        string `json:"bank"`
}

type ProductInventory struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
