// routes/router.go

package routes // routes: package for all routes

import (
	"myfiberproject/handlers"   // handlers: package for request handlers
	"myfiberproject/middleware" // middleware: package for middleware

	"github.com/gofiber/fiber/v2" // fiber: package for building web applications
)

const (
	BaseUserPath = "/users"
	UserByIDPath = "/users/:id"
)

func SetupRoutes(app *fiber.App) { // SetupRoutes: function to set up all routes

	// User routes
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)
	app.Post("/forgot-password", handlers.ForgotPassword)

	// Admin routes
	app.Post("/seed-admin", handlers.SeedAdminHandler)

	// Users management route
	app.Get(BaseUserPath, middleware.RequireRole([]string{"administrator"}, "approved"), handlers.GetAllUsers)
	app.Get(UserByIDPath, middleware.RequireRole([]string{"administrator"}, "approved"), handlers.GetSingleUser)
	app.Put(UserByIDPath, middleware.RequireRole([]string{"administrator"}, "approved"), handlers.UpdateUser)
	app.Delete(UserByIDPath, middleware.RequireRole([]string{"administrator"}, "approved"), handlers.DeleteUser)

	// Article category
	app.Post("/article-category", middleware.RequireRole([]string{"administrator"}, "approved"), handlers.CreateArticleCategory)
	app.Get("/article-category", middleware.RequireRole([]string{"administrator"}, "approved"), handlers.GetAllArticleCategory)

}
