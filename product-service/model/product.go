package model

type Product struct {
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
	Image       string  `bson:"image"`
}

