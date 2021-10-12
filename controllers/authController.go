package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/amshashankk/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

var client *mongo.Client

func Register(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	var data map[string]string
	var user models.User

	//fmt.Println("Check1")

	if err := c.BodyParser(data); err != nil {
		fmt.Println("Error Occured here")
	}

	//fmt.Println("Check2")

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 10)

	user = models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	collection := client.Database("Shashank").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, &user)
	if result != nil {
		fmt.Println("Checking")
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	var data map[string]string
	var dbUser models.User
	var user models.User

	if err := c.BodyParser(data); err != nil {
		return err
	}

	collection := client.Database("Shashank").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return err
	}

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User Not Registered",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect Password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Authentication Failed",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	//database.DB.Where("id = ?", claims.Issuer).(&user)
	collection := client.Database("Shashank").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = collection.FindOne(ctx, claims.Id).Decode(&user)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
