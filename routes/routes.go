package routes

import (
	"Ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoutes general routes
func UserRoutes(request *gin.Engine) {
	request.POST("/users/signup", controllers.SignUp())
	request.POST("/users/login", controllers.Login())
	request.POST("/admin/add-product", controllers.ProductViewerAdmin())
	request.GET("/users/product-view", controllers.SearchProduct())
	request.POST("/users/search", controllers.SearchProductByQuery())
}
