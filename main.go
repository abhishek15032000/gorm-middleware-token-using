package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var jwtSecret = []byte("nsdkjnsakndsakndsakndsakdnsak")

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json : "name"`
	Email    string `json : "email" gorm:"unique"`
	Password string `json : "-"` // hide in response
}

func main() {
	dsn := "host=localhost user=postgres password=1234 dbname=testdb port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db = database
	// migrate the user model
	db.AutoMigrate(&User{})
	e := echo.New()
	// middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.POST("/register", registerUser)
	e.POST("/login", loginUser)
	r := e.Group("/user")
	r.Use(authMiddleware)
	r.GET("/profile", profile)
	e.Logger.Fatal(e.Start(":8091"))
}

func registerUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// password hashing // bcrypt used for hashing

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to hash the password",
		})
	}
	u.Password = string(hash)

	if err := db.Create(&u).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Email already exists",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{"error": "user registered succesfully"})
}

func loginUser(c echo.Context) error {
	req := new(User)
	if err := c.Bind(req); err != nil {
		return err
	}
	// decryption of hash
	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid Email or Password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid Email or Password",
		})
	}

	// token generate

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to generate token",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Loggin successfull",
		"token":   t,
	})
}

func profile(c echo.Context) error {
	userId := c.Get("user_id")
	// fmt.Fprintf(c.Response(), "%v", userId)
	var user User
	if err := db.First(&user, userId).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}
	return c.JSON(http.StatusOK, user)
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "missing token",
			})
		}
		tokenString := ""
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "invalid token",
			})
		}
		return next(c)
	}
}
