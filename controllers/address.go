package controllers

import (
	"Ecommerce/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"time"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId = c.Query("id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}
		address, err := ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var addresses models.Address
		addresses.AddressId = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unWind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointerCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unWind, grouping})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var addressInfo []bson.M
		if err = pointerCursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}
		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err = UserCollection.UpdateOne(ctx, filter, update)

			if err != nil {
				fmt.Println(err)
			}

		} else {
			c.IndentedJSON(http.StatusBadRequest, "Not allowed to add less than 2 addresses")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId = c.Query("id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{key: "_id", Value: userObjectId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street},
			{Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.PinCode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Home address updated successfully")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId = c.Query("id")
		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{key: "_id", Value: userObjectId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street},
			{Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.PinCode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "Work address updated successfully")
	}
}

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
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
		}
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userObjectId}}
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
