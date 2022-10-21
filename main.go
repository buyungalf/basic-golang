package main

import (
	"rapidtech/shoppingcart/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
)

// var checker = validator.New()

func main()  {
	store := session.New()

	engine := html.New("./views",".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/","./public")

	prodController := controllers.InitProductController(store)
	authController := controllers.InitAuthController(store)
	cartController := controllers.InitCartController(store)

	prod := app.Group("/products")
	prod.Get("/", prodController.IndexProducts)
	prod.Get("/detail/:id", prodController.DetailProduct)
	prod.Get("/create", prodController.AddProduct)
	prod.Post("/create", prodController.AddPostedProduct)
	prod.Get("/edit/:id", prodController.EditProduct)
	prod.Post("/edit/:id", prodController.EditPostedProduct)
	prod.Get("/delete/:id", prodController.DeleteProduct)

	app.Get("/register", authController.Register)
	app.Get("/profile", authController.Profile)
	app.Post("/register", authController.PostRegister)
	app.Get("/check", authController.CheckSession)
	app.Get("/logout", authController.Logout)

	app.Get("/cart", cartController.GetCart)
	app.Get("/addtocart", cartController.AddtoCart)

	log := app.Group("login")
	log.Get("/", authController.Login)
	log.Post("/", authController.PostLogin)

	app.Listen(":3000")
}