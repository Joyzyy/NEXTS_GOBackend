package controllers

import (
	"context"
	"example/hello/configs"
	"example/hello/models"
	"example/hello/responses"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection *mongo.Collection = configs.GetDB(configs.CLIENT, "product", "productData")

type HelperFunctions interface {
	GetProduct(product []models.Product)
	InternalServerError(err error)
	BadRequest(err error)
}

type ProductControllers struct {
	HelperFunctions
	g *gin.Context
}

func (p *ProductControllers) GetProduct(product []models.Product) {
	p.g.JSON(
		http.StatusOK,
		responses.ProductResponse{
			Status:  http.StatusOK,
			Message: "",
			Data:    map[string]interface{}{"data": product},
		},
	)
}

func (p *ProductControllers) InternalServerError(err error) {
	p.g.JSON(
		http.StatusInternalServerError,
		responses.ProductResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Data:    map[string]interface{}{"error": err},
		},
	)
}

func (p *ProductControllers) BadRequest(err error) {
	p.g.JSON(
		http.StatusBadRequest,
		responses.ProductResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Data:    map[string]interface{}{"error": err},
		},
	)
}

func CreateProduct() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		p := ProductControllers{g: g}
		var product models.Product
		defer cancel()

		if err := g.BindJSON(&product); err != nil {
			p.BadRequest(err)
			return
		}

		if validationError := validator.New().Struct(product); validationError != nil {
			p.BadRequest(validationError)
			return
		}

		newProduct := models.Product{
			Id:          primitive.NewObjectID(),
			Name:        product.Name,
			TabType:     product.TabType,
			Image:       product.Image,
			Category:    product.Category,
			Sizes:       product.Sizes,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		}

		_, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			p.InternalServerError(err)
			return
		}

		p.GetProduct([]models.Product{newProduct})
	}
}

func GetAllProductsByFilters(g *gin.Context, p ProductControllers, ctx context.Context, products []models.Product, queries map[string]interface{}) {
	filter := bson.M{}

	var cursor *mongo.Cursor
	var err error
	var opts *options.FindOptions

	for key, value := range queries {
		if value != "" {
			filter[key] = value

			if key == "price" {
				price, _ := strconv.Atoi(value.(string))
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
						filter[key] = bson.M{"$gte": price}
					} else if priceTypeQuery == "lte" {
						filter[key] = bson.M{"$lte": price}
					} //@TODO: range query
				} else {
					filter[key] = price
				}
			}

			if key == "sizes" {
				sizeQueryArray := strings.Split(value.(string), ",")
				sizeArray := make([]int, 0)

				for _, size := range sizeQueryArray {
					size, _ := strconv.Atoi(size)
					sizeArray = append(sizeArray, size)
				}

				filter[key] = bson.M{"$in": sizeArray}
			}
		}
	}

	cursor, err = productCollection.Find(ctx, filter, opts)
	if err != nil {
		p.InternalServerError(err)
		return
	}

	if err = cursor.All(ctx, &products); err != nil {
		p.InternalServerError(err)
		return
	}

	p.GetProduct(products)
}

func GetAllProducts() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		p := ProductControllers{g: g}
		var products []models.Product
		defer cancel()

		// check if there is a query
		if g.Query("category") != "" || g.Query("price") != "" || g.Query("sizes") != "" {
			GetAllProductsByFilters(g, p, ctx, products,
				map[string]interface{}{
					"category": g.Query("category"),
					"price":    g.Query("price"),
					"sizes":    g.Query("sizes"),
				},
			)
			return
		}

		cursor, err := productCollection.Find(ctx, models.Product{})
		if err != nil {
			p.InternalServerError(err)
			return
		}

		if err = cursor.All(ctx, &products); err != nil {
			p.InternalServerError(err)
			return
		}

		p.GetProduct(products)
	}
}

func FindProductByName() gin.HandlerFunc {
	return func(g *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productName := g.Param("name")
		p := ProductControllers{g: g}
		var product models.Product
		defer cancel()

		err := productCollection.FindOne(ctx, bson.M{"name": productName}).Decode(&product)
		if err != nil {
			p.InternalServerError(err)
			return
		}

		p.GetProduct([]models.Product{product})
	}
}
