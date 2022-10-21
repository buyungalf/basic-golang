package controllers

import (
	"rapidtech/shoppingcart/database"
	"rapidtech/shoppingcart/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	Db *gorm.DB
	store *session.Store
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Register",
	})
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func (controller *AuthController) PostLogin(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)

	if err!=nil {
		panic(err)
	}

	var myform models.User
	var data models.User

	if err := c.BodyParser(&myform); err != nil {
		return c.JSON(fiber.Map{"error": err})
	}

	username := myform.Username
	plainPassword := myform.Password

	err2 := models.ReadOneUser(controller.Db, &data, username)

	if err2 != nil {
		return c.Redirect("/login")
	}
	
	hashPassword := data.Password

	check := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))

	status := check == nil

	if status {
		sess.Set("username", username)
		sess.Save()
		return c.Redirect("/products")
	} else {
		return c.Redirect("/login")
	}

	
}

func (controller *AuthController) PostRegister(c *fiber.Ctx) error {
	var register models.User

		if err := c.BodyParser(&register); err != nil {
			return c.Redirect("/register")
		}

		bytes, _ := bcrypt.GenerateFromPassword([]byte(register.Password), 8)
		sHash := string(bytes)
		
		register.Password = sHash

		err := models.Register(controller.Db, &register)

		if err != nil {
			return c.Redirect("/register")
		}
		
		return c.Redirect("/login")
}

// /profile
func (controller *AuthController) Profile(c *fiber.Ctx) error {

	var users []models.User
	err := models.ReadUser(controller.Db, &users)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"Title": "Users",
		"Products": users,
	})
}
// /logout
func (controller *AuthController) Logout(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)

	if err != nil {
		panic(err)
	}

	sess.Destroy()

	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func InitAuthController(s *session.Store) *AuthController {
	db := database.InitDb()

	db.AutoMigrate(&models.User{})
	
	return &AuthController{Db: db, store: s}
}