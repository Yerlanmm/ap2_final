package handler

import (
	"context"
	"html/template"
	"net/http"
	"time"
	"product-service/config"
	"product-service/model"

	"go.mongodb.org/mongo-driver/bson"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Collection("product")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to load products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		http.Error(w, "Failed to parse products", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/main.html"))
	err = tmpl.Execute(w, products)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}