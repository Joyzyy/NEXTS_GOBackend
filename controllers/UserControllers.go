package controllers

import (
	"context"
	"example/hello/configs"
	"example/hello/models"
	"example/hello/responses"
	"example/hello/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetDB(configs.CLIENT, "user", "users")

func Register() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := g.BindJSON(&user); err != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "No content-type/json provided.",
					Data:    map[string]interface{}{"data": ""},
				},
			)
			return
		}

		if validationError := validator.New().Struct(user); validationError != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Validation Error",
					Data:    map[string]interface{}{"data": validationError.Error()},
				},
			)
			return
		}

		usernameCount, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			log.Fatal(err)
			return
		}
		if usernameCount > 0 {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Username already exists.",
					Data:    map[string]interface{}{"data": ""},
				},
			)
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal(err)
			return
		}
		if emailCount > 0 {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Email already exists.",
					Data:    map[string]interface{}{"data": ""},
				},
			)
			return
		}

		uid := primitive.NewObjectID()

		token, refreshToken, err := utils.Issue(uid.Hex())
		if err != nil {
			log.Fatal(err)
			return
		}

		newUser := models.User{
			Id:         uid,
			Username:   user.Username,
			Password:   utils.HashPassword(user.Password),
			Email:      user.Email,
			Created_at: time.Now().Local(),
		}

		res, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			g.JSON(
				http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		g.SetCookie("refreshToken", refreshToken, 60*60*24*7, "/", "localhost", false, true)
		g.JSON(
			http.StatusCreated,
			responses.UserResponse{
				Status:  http.StatusCreated,
				Message: "Created",
				Data:    map[string]interface{}{"jwt": token, "data": res.InsertedID},
			},
		)
	}
}

func Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var userDB models.User
		defer cancel()

		if err := g.BindJSON(&user); err != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Bad request",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userDB)
		if err != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Username or password are incorrect",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		isPasswordValid, message := utils.VerifyPassword(user.Password, userDB.Password)
		defer cancel()
		if !isPasswordValid {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Username or password are incorrect",
					Data:    map[string]interface{}{"data": message},
				},
			)
			return
		}
		if userDB.Username == "" {
			g.JSON(
				http.StatusBadRequest,
				responses.UserResponse{
					Status:  http.StatusBadRequest,
					Message: "Username or password are incorrect",
					Data:    map[string]interface{}{"data": ""},
				},
			)
			return
		}

		token, refreshToken, err := utils.Issue(userDB.Id.Hex())
		if err != nil {
			log.Fatal(err)
			return
		}

		g.SetSameSite(http.SameSiteNoneMode)
		g.SetCookie("token", token, 60*60*24, "/", "https://nextjs-app-charka-frontend.herokuapp.com", true, false)
		g.SetCookie("refreshToken", refreshToken, 60*60*24*7, "/", "https://nextjs-app-charka-frontend.herokuapp.com", true, true)
		fmt.Println("asd")
		g.JSON(
			http.StatusOK,
			responses.UserResponse{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    map[string]interface{}{"data": token},
			},
		)
	}
}

func Logout() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.SetSameSite(http.SameSiteNoneMode)
		g.SetCookie("refreshToken", "data", 0, "/", "https://nextjs-app-charka-frontend.herokuapp.com", true, true)
		g.JSON(
			http.StatusOK,
			responses.UserResponse{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    map[string]interface{}{"data": "data"},
			},
		)
	}
}
