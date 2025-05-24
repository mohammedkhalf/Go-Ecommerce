package main

import (
	"Ecommerce/controllers"
	"Ecommerce/database"
	"Ecommerce/middleware"
	"Ecommerce/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Product"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// routes need auth
	router.GET("/add-to-cart", app.AddToCart())
	router.GET("/remove-item", app.RemoveItem())
	router.GET("/cart-checkout", app.BuyFromCart())
	router.GET("/instance-buy", app.InstanceBuy())

	log.Fatal(router.Run(":" + port))

}
