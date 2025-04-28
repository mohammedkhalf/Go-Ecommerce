package routes

import (
	"Ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

// general routes
func UserRoutes(request *gin.Engine) {
	request.POST("/user/signup", controllers.SignUp())
	request.POST("/users/login", controllers.Login())
	request.POST("/admin/add-product", controllers.AddProduct())
	request.GET("/users/product-view", controllers.ViewProduct())
	request.POST("/users/search", controllers.SearchProduct())
}
