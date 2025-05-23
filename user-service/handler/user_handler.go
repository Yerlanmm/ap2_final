package handler

import (
	"context"
	"html/template"
	"net/http"
	"time"
	"user-service/config"
	"user-service/model"
	"user-service/utils"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}


func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear cookies or session logic if needed
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		http.Error(w, "Please fill all fields", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing model.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)
	if err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8001/main", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Please fill all fields", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser model.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	_, err = utils.GenerateJWT(dbUser.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8001/main", http.StatusSeeOther)
}
