package controllers

import (
	"Ecommerce/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"net/http"
	"time"
)

func HashPassword(password string) string {}

func VerifyPassword(userPassword string, givenPassword string) (bool, error) {}

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJson(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
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

		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.newObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)

		_, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this user didn't created "})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
	}

}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.withTimeout(context.Background(), 5*time.Second)
		defer cacnel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.findOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		isPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		defer cancel()

		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refershToken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId)

		defer cancel()
		generate.UpdateAllToken(token, refershToken, foundUser.UserId)
		c.JSON(http.StatusFound, foundUser)
	}
}

func AddProduct() gin.HandlerFunc {
}

func ViewProduct() gin.HandlerFunc {
}

func SearchProduct() gin.HandlerFunc {
}
