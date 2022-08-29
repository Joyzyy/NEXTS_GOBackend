package controllers

import (
	"context"
	"example/hello/configs"
	"example/hello/models"
	"example/hello/responses"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			TabType:     product.TabType,
			Image:       product.Image,
			Category:    product.Category,
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

func GetAllProductsByCategory(g *gin.Context, ctx context.Context, products []models.Product, categoryQuery string) {
	cursor, err := productCollection.Find(ctx, bson.M{"category": categoryQuery})
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

func GetAllProductsByPrice(g *gin.Context, ctx context.Context, products []models.Product, priceQuery string) {
	price, _ := strconv.Atoi(priceQuery)

	filter := bson.M{"price": price}
	var opts *options.FindOptions

	if optsQuery := g.Query("ord"); optsQuery != "" {
		if optsQuery == "asc" {
			opts = options.Find().SetSort(bson.M{"price": 1})
		} else if optsQuery == "desc" {
			opts = options.Find().SetSort(bson.M{"price": -1})
		}
	} else {
		opts = nil
	}

	if priceTypeQuery := g.Query("type"); priceTypeQuery != "" {
		if priceTypeQuery == "gte" {
			filter = bson.M{"price": bson.M{"$gte": price}}
		} else {
			filter = bson.M{"price": bson.M{"$lte": price}}
		}
	}

	cursor, err := productCollection.Find(ctx, filter, opts)
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

func GetAllProductsByCategoryAndPrice(g *gin.Context, ctx context.Context, products []models.Product, categoryQuery string, priceQuery string) {
	price, _ := strconv.Atoi(priceQuery)

	var cursor *mongo.Cursor
	var err error
	var filter bson.M
	var opts *options.FindOptions

	if optsQuery := g.Query("ord"); optsQuery != "" {
		if optsQuery == "asc" {
			opts = options.Find().SetSort(bson.M{"price": 1})
		} else if optsQuery == "desc" {
			opts = options.Find().SetSort(bson.M{"price": -1})
		}
	} else {
		opts = nil
	}

	filter = bson.M{"price": price, "category": categoryQuery}

	if priceTypeQuery := g.Query("type"); priceTypeQuery != "" {
		if priceTypeQuery == "gte" {
			filter = bson.M{"price": bson.M{"$gte": price}, "category": categoryQuery}
		} else {
			filter = bson.M{"price": bson.M{"$lte": price}, "category": categoryQuery}
		}
	}

	cursor, err = productCollection.Find(ctx, filter, opts)
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

func GetAllProducts() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var products []models.Product
		categoryQuery := g.Query("category")
		priceQuery := g.Query("price")
		defer cancel()

		if categoryQuery != "" && priceQuery == "" {
			GetAllProductsByCategory(g, ctx, products, categoryQuery)
			return
		} else if categoryQuery != "" && priceQuery != "" {
			GetAllProductsByCategoryAndPrice(g, ctx, products, categoryQuery, priceQuery)
			return
		} else if categoryQuery == "" && priceQuery != "" {
			GetAllProductsByPrice(g, ctx, products, priceQuery)
			return
		}

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

func FindProductByName() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		productName := g.Param("name")
		var product models.Product
		defer cancel()

		err := productCollection.FindOne(ctx, bson.M{"name": productName}).Decode(&product)
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
