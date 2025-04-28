package controllers

import (
	"github.com/gin-gonic/gin"
)

func HashPassword(password string) string {
}

func VerifyPassword(userPassword string, givenPassword string) (bool, error) {
}

func SignUp() gin.HandlerFunc {
}

func Login() gin.HandlerFunc {
}

func AddProduct() gin.HandlerFunc {
}

func ViewProduct() gin.HandlerFunc {
}

func SearchProduct() gin.HandlerFunc {
}
