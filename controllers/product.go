package controllers

import (
	"rapidtech/shoppingcart/database"
	"rapidtech/shoppingcart/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type ProductController struct {
	Db *gorm.DB
}

//GET
func (controller *ProductController) IndexProducts(c *fiber.Ctx) error {
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.Render("product", fiber.Map{
		"Title": "Products",
		"Products": products,
	})
}

func (controller *ProductController) DetailProduct(c *fiber.Ctx) error {
	var product models.Product
	var id, err = c.ParamsInt("id")
	
	if err != nil {
		return err
	}

	err2 := models.ReadProductById(controller.Db, &product, id)

	if err2 != nil {
		return c.SendStatus(500)
	}

	return c.Render("productDetail", fiber.Map{
		"Title": "Products",
		"Product": product,
	})
}

func (controller *ProductController) AddProduct(c *fiber.Ctx) error {
	return c.Render("addProduct", fiber.Map{
		"title": "Tambah Product",
	})
}

func (controller *ProductController) EditProduct(c *fiber.Ctx) error {
	var product models.Product
	var id, err = c.ParamsInt("id")
	
	if err != nil {
		return err
	}

	err2 := models.ReadProductById(controller.Db, &product, id)

	if err2 != nil {
		return c.SendStatus(500)
	}

	return c.Render("editProduct", fiber.Map{
		"Title": "Products",
		"Product": product,
	})
}

func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	var product models.Product
	var id, err = c.ParamsInt("id")
	
	if err != nil {
		return err
	}

	err2 := models.DeleteProduct(controller.Db, &product, id)

	if err2 != nil {
		return c.SendStatus(500)
	}

	return c.Redirect("/products")
}


//POST
func (controller *ProductController) AddPostedProduct(c *fiber.Ctx) error {
	// data := new(models.Product)
	var data models.Product

		if err := c.BodyParser(&data); err != nil {
			return c.Redirect("/products")
		}

		err := models.CreateProduct(controller.Db, &data)

		if err != nil {
			return c.Redirect("/products")
		}

		return c.Redirect("/products")
}

func (controller *ProductController) EditPostedProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)


	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var myform models.Product

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}
	product.Name = myform.Name
	product.Quantity = myform.Quantity
	product.Price = myform.Price
	// save product
	models.UpdateProduct(controller.Db, &product)
	
	return c.Redirect("/products")	
}





func InitProductController() *ProductController {
	db := database.InitDb()

	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db}
}