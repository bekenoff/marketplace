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
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Cart struct {
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

type Sight struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	ContentInfo  string `json:"content_info"`
	BusNumbers   string `json:"bus_numbers"`
	WorkingDays  string `json:"working_days"`
	WorkingHours string `json:"working_hours"`
	Visited      int    `json:"visited"`
	ImageUrl     string `json:"image_url"`
}
type Client struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	First_name string `json:"firstname"`
	Last_name  string `json:"lastname"`
	Telephone  string `json:"telephone"`
}

type Client_Sights struct {
	Id       int `json:"id"`
	ClientId int `json:"client_id"`
	SightId  int `json:"sight_id"`
}

type Client_Events struct {
	Id       int `json:"id"`
	ClientId int `json:"client_id"`
	EventId  int `json:"event_id"`
}

type Events struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Schedule    string `json:"schedule"`
	Bus         string `json:"bus"`
	Price       string `json:"price"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Category    int    `json:"category"`
}

type EventsCategory struct {
	Id            int `json:"id"`
	EventCategory int `json:"event_category"`
}

type Recommendation struct {
	Id              int `json:"id"`
	ClientId        int `json:"client_id"`
	SightCategoryId int `json:"sight_category_id"`
	EventCategoryId int `json:"event_category_id"`
}

type SightsCategory struct {
	Id            int `json:"id"`
	SightCategory int `json:"sight_category"`
}
