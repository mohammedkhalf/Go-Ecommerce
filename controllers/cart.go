package controllers

import (
	"Ecommerce/database"
	"Ecommerce/models"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"time"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func newApplication(productCollection *mongo.Collection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.Handler {
	return func(c *gin.Context) {
		var productQueryID = c.Query("productID")
		if productQueryID == "" {
			log.Println("productID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is empty"))
			return
		}
		var userQueryID = c.Query("userID")
		if userQueryID == "" {
			log.Println("userID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("productID is invalid"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.AddProductToCart(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "product added to cart")
	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productQueryID = c.Query("productID")
		if productQueryID == "" {
			log.Println("productID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is empty"))
			return
		}
		var userQueryID = c.Query("userID")
		if userQueryID == "" {
			log.Println("userID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("productID is invalid"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "product removed from cart")
	}

}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Id"})
			c.Abort()
			return
		}
		userObjectID, _ = primitive.ObjectIDFromHex(userId)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userObjectID}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, "Item not found")
			return
		}
		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userObjectID}}}}
		unWind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}}}}}

		pointerCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unWind, grouping})
		if err != nil {
			log.Println(err)
		}
		var Listing []bson.M
		if err = pointerCursor.All(ctx, &Listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range Listing {
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userQueryID = c.Query("userID")
		if userQueryID == "" {
			log.Println("userID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully Placed the order")
	}
}

func (app *Application) InstanceBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productQueryID = c.Query("productID")
		if productQueryID == "" {
			log.Println("productID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is empty"))
			return
		}
		var userQueryID = c.Query("userID")
		if userQueryID == "" {
			log.Println("userID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("productID is invalid"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.InstanceBuyer(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully Placed the order")
	}
}
