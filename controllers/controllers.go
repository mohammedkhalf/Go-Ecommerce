package controllers

import (
	"Ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"net/http"
	"time"
)

func HashPassword(password string) string {
}

func VerifyPassword(userPassword string, givenPassword string) (bool, error) {
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}

		validationErr := Validate.Struct(&user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.countDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}

		count, err = UserCollection.countDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number already exists"})
			return
		}

	}

}

func Login() gin.HandlerFunc {
}

func AddProduct() gin.HandlerFunc {
}

func ViewProduct() gin.HandlerFunc {
}

func SearchProduct() gin.HandlerFunc {
}
