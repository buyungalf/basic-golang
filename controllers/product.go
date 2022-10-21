package controllers

import (
	"fmt"
	"rapidtech/shoppingcart/database"
	"rapidtech/shoppingcart/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)


type ProductController struct {
	Db *gorm.DB
	store *session.Store
}

//GET
func (controller *ProductController) IndexProducts(c *fiber.Ctx) error {
	
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)

	if err != nil {
		return c.SendStatus(500)
	}

	sess,_ := controller.store.Get(c)
	id := sess.Get("id")

	return c.Render("product", fiber.Map{
		"Title": "Products",
		"Products": products,
		"UserId": id,
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
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]
		
		for _, file := range files {
			var data models.Product
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			
			if err := c.BodyParser(&data); err != nil {
				return c.Redirect("/products")
			}
			
			if err := c.SaveFile(file, fmt.Sprintf("./public/upload/%s", file.Filename)); err != nil {
				return err
			}

			data.Image = file.Filename
		
			err := models.CreateProduct(controller.Db, &data)
		
			if err != nil {
				return c.Redirect("/products")
			}

			c.Redirect("/products")
		}
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "error",
	})
}


func (controller *ProductController) EditPostedProduct(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]
		
		for _, file := range files {
			var data models.Product
			
			id := c.Params("id")
			idn,_ := strconv.Atoi(id)


			err := models.ReadProductById(controller.Db, &data, idn)
			if err!=nil {
				return c.SendStatus(500) // http 500 internal server error
			}
			
			var myform models.Product

			if err := c.BodyParser(&myform); err != nil {
				return c.Redirect("/products")
			}

			data.Name = myform.Name
			data.Quantity = myform.Quantity
			data.Price = myform.Price
			// save product
			

			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			
			if err := c.SaveFile(file, fmt.Sprintf("./public/upload/%s", file.Filename)); err != nil {
				return err
			}

			data.Image = file.Filename
		
			err = models.UpdateProduct(controller.Db, &data)
		
			if err != nil {
				// return c.Redirect("/products")
				return c.JSON(data)
			}

			c.Redirect("/products")
		}
		return c.JSON(fiber.Map{
			"message": "error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "error",
	})
		
	// return c.Redirect("/products")	
}





func InitProductController(s *session.Store) *ProductController {
	db := database.InitDb()

	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db, store: s}
}