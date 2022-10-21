package controllers

import (
	"rapidtech/shoppingcart/database"
	"rapidtech/shoppingcart/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)


type CartController struct {
	Db *gorm.DB
	store *session.Store
}

//GET
func (controller *CartController) GetCart(c *fiber.Ctx) error {
	
	var cart []models.Cart
	err := models.ViewCart(controller.Db, &cart)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.Render("cart", fiber.Map{
		"Title": "Cart",
		"Cart": cart,
	})
}

func (controller *CartController) AddtoCart(c *fiber.Ctx) error {
	userId := c.Query("userid")
	productId := c.Query("productid")

	product_id,_ := strconv.Atoi(productId)
	user_id, err := strconv.Atoi(userId)
	
	if err != nil {
		c.Redirect("/login")
	}
  
	var cart models.Cart
	var product models.Product

	err2 := models.ReadProductById(controller.Db, &product, product_id)

	if err2 != nil {
		c.Redirect("/products")
	}

	cart.ProductId = product_id
	cart.UserId = user_id
	cart.Quantity = 1
	cart.Total = float64(cart.Quantity)*product.Price

	err3 := models.AddtoCart(controller.Db, &cart)

	if err3 != nil {
		c.Redirect("/products")
	}

	return c.Redirect("/cart")
}

func InitCartController(s *session.Store) *CartController {
	db := database.InitDb()

	db.AutoMigrate(&models.Cart{})

	return &CartController{Db: db, store: s}
}