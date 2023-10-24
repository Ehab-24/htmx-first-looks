package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Ehab-24/test/db"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	PasswordHash []byte             `bson:"password_hash"`
}

type LoginPayload struct {
	Email    string
	Password string
}

type RegisterPayload struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

func getAuthRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", login)
	r.Get("/user", user)
	r.Post("/register", registerUser)

	// TODO: remove this route
	r.Get("/users", getUsers)
	return r
}

//***************** Handlers *****************

func getUsers(w http.ResponseWriter, r *http.Request) {
	filter := bson.D{}
	cursor, err := db.Db.Collection("users").Find(context.Background(), filter)
	if CheckAPIError(w, err) {
		return
	}
	var users []User
	cursor.All(context.Background(), &users)
	json.NewEncoder(w).Encode(users)
}

func login(w http.ResponseWriter, r *http.Request) {
	payload := LoginPayload{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// validate payload
	if payload.Password == "" || payload.Email == "" {
		w.Write([]byte("missing required field(s)"))
		return
	}

	// verify credentials
	filter := bson.D{{Key: "email", Value: payload.Email}}
	_user := db.Db.Collection("users").FindOne(context.Background(), filter)
	var user User
	_user.Decode(&user)
	if user.Email == "" {
		w.Write([]byte("user doesn't exist"))
		return
	}
	hashedPass := sha256.New().Sum([]byte(payload.Password))
	if string(hashedPass) != string(user.PasswordHash) {
		w.Write([]byte("wrong password"))
		return
	}

	// generate access token
	accessToken, err := generateJWT(os.Getenv("JWT_SECRET"), "user_id", user.ID.Hex())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// set cookie on client browser
	cookie := http.Cookie{
		Name:     "x-access-token",
		Value:    accessToken,
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteStrictMode,
		Secure:   os.Getenv("ENV") != "development",
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte("success"))
}

func user(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("x-access-token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("access token: " + cookie.Value))
}

// register user
func registerUser(w http.ResponseWriter, r *http.Request) {
	payload := RegisterPayload{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	// validate payload
	if payload.Password == "" || payload.Email == "" || payload.Name == "" {
		w.Write([]byte("missing required field(s)"))
		return
	}
	if payload.Password != payload.ConfirmPassword {
		w.Write([]byte("passwords don't match"))
		return
	}

	// check for duplicate
	filter := bson.D{{Key: "email", Value: payload.Email}}
	_duplicate := db.Db.Collection("users").FindOne(context.Background(), filter)
	var duplicate User
	_duplicate.Decode(&duplicate)
	if duplicate.Email == payload.Email {
		w.Write([]byte("user already exists"))
		return
	}

	// create new user
	hashedPass := sha256.New().Sum([]byte(payload.Password))
	newUser := NewUser(payload.Name, payload.Email, hashedPass)
	_, err := db.Db.Collection("users").InsertOne(context.Background(), newUser)
	if CheckAPIError(w, err) {
		return
	}

	// geenrate access token
	accessToken, err := generateJWT(os.Getenv("JWT_SECRET"), "user_id", newUser.ID.Hex())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// set cookie on client browser
	cookie := http.Cookie{
		Name:     "x-access-token",
		Value:    accessToken,
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteStrictMode,
		Secure:   os.Getenv("ENV") != "development",
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte("success"))
}

//***************** Helpers *****************

func getClaims(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	json.NewEncoder(w).Encode(claims)
}

func NewUser(name string, email string, passwordHash []byte) User {
	return User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
	}
}
