package controllers

import (
	"Ecommerce/models"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"net/http"
	"time"
)

func AddAddress() gin.HandlerFunc {}

func EditHomeAddress() gin.HandlerFunc {}

func EditWorkAddress() gin.HandlerFunc {}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId = c.Query("id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Search Index"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		userId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "address deleted successfully")
	}

}
