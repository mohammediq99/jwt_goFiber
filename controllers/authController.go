package controllers

import (
	"project1/database"
	"project1/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secretkey"

func Hello(c *fiber.Ctx) error {
	return c.SendString("Helloe")
}
func Register(c *fiber.Ctx) error {
	var data map[string]string

	if error := c.BodyParser(&data); error != nil {
		return error
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	database.DB.Create(&user)
	return c.JSON(user)

	// return c.SendString("Helloe")/
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if error := c.BodyParser(&data); error != nil {
		return error
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if error := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Password mismatch",
		})

	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, error := claims.SignedString([]byte(SecretKey))

	if error != nil {

		if error := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "couldn't Login",
			})

		}

	}

	cookie := fiber.Cookie{
		Name:     "JWT",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successfully LogedIn",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("JWT")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "not Authorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "JWT",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(cookie)
}
