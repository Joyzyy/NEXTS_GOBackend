package controllers

import (
	"context"
	"example/hello/configs"
	"example/hello/models"
	"example/hello/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetDB(configs.CLIENT, "product", "productData")

func CreateProduct() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var product models.Product
		defer cancel()

		if err := g.BindJSON(&product); err != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.ProductResponse{
					Status:  http.StatusBadRequest,
					Message: "Bad Request",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		if validationError := validator.New().Struct(product); validationError != nil {
			g.JSON(
				http.StatusBadRequest,
				responses.ProductResponse{
					Status:  http.StatusBadRequest,
					Message: "Bad Request",
					Data:    map[string]interface{}{"data": validationError.Error()},
				},
			)
			return
		}

		newProduct := models.Product{
			Id:          primitive.NewObjectID(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		}

		res, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			g.JSON(
				http.StatusInternalServerError,
				responses.ProductResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		g.JSON(
			http.StatusCreated,
			responses.ProductResponse{
				Status:  http.StatusCreated,
				Message: "Created",
				Data:    map[string]interface{}{"data": res.InsertedID},
			},
		)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var products []models.Product
		defer cancel()

		cursor, err := productCollection.Find(ctx, models.Product{})
		if err != nil {
			g.JSON(
				http.StatusInternalServerError,
				responses.ProductResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		if err = cursor.All(ctx, &products); err != nil {
			g.JSON(
				http.StatusInternalServerError,
				responses.ProductResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		g.JSON(
			http.StatusOK,
			responses.ProductResponse{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    map[string]interface{}{"data": products},
			},
		)
	}
}

func FindProductByID() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		productID := g.Param("id")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productID)

		err := productCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&product)
		if err != nil {
			g.JSON(
				http.StatusInternalServerError,
				responses.ProductResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		g.JSON(
			http.StatusOK,
			responses.ProductResponse{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    map[string]interface{}{"data": product},
			},
		)
	}
}
