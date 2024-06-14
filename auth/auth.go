package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"tarjeta/jwt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LoginDetails struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

type User struct {
	UUID           string `json:"uuid" xml:"uuid" form:"uuid"`
	Username       string `json:"username" xml:"username" form:"username"`
	Email          string `json:"email" xml:"email" form:"email"`
	Role           string `json:"role" xml:"role" form:"role"`
	ProfilePicture string `json:"profilePicture" xml:"profilePicture" form:"profilePicture"`
}

func constructCookie(key string, value string, secure bool) (cookie *fiber.Cookie) {
	cookie = new(fiber.Cookie)
	cookie.Name = key
	if value != "" {
		cookie.Value = value
	}
	cookie.SessionOnly = true
	cookie.Path = "/"
	cookie.Secure = secure
	cookie.HTTPOnly = false
	cookie.SameSite = "lax"
	return cookie
}

// POST
func Login(ctx *fiber.Ctx) error {
	p := new(LoginDetails)
	if err := ctx.BodyParser(p); err != nil {
		log.Print(err)
		return ctx.Status(http.StatusBadRequest).SendString("An error occured!")
	}
	// TODO: Add database
	if p.Username != "alex" && p.Password != "qwerty123$!" {
		return ctx.Status(http.StatusUnauthorized).SendString("Wrong credentials")
	}
	token, err := jwt.EncryptJWT(p.Username)
	if err != nil {
		log.Print(err)
		return ctx.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}
	session := constructCookie("jwt", token.String(), true)
	ctx.Cookie(session)
	return ctx.Redirect("/profile.html", http.StatusMovedPermanently)
}

// POST
func Logout(ctx *fiber.Ctx) error {
	session := constructCookie("jwt", "", true)
	session.SessionOnly = false
	session.Expires = time.Now().Add(-(time.Hour * 2))
	user := constructCookie("user", "", false)
	user.SessionOnly = false
	user.Expires = time.Now().Add(-(time.Hour * 2))
	ctx.Cookie(session)
	ctx.Cookie(user)
	return ctx.Redirect("/login.html", http.StatusMovedPermanently)
}

func WhoAmI(ctx *fiber.Ctx) error {
	seed := ctx.Cookies("jwt")
	if seed == "" {
		return ctx.Status(http.StatusUnauthorized).SendString("You are not authorized!")
	}
	token, err := jwt.DecryptJWT(seed)
	if err != nil {
		log.Printf("[Error] WhoAmI | %v\n", err)
		return ctx.Status(http.StatusBadRequest).SendString("An error occured!")
	}
	if token.IsExpired() {
		ctx.ClearCookie("jwt")
		return ctx.Status(http.StatusUnauthorized).SendString("You are not authorized!")
	}
	user := &User{uuid.NewString(), "Alex W. Gaming", "me@alexalex.dev", "Admin", "https://avatars.githubusercontent.com/u/21964985?v=4"}
	userByte, err := json.Marshal(user)
	if err != nil {
		log.Printf("[Error] WhoAmI | %v\n", err)
		return ctx.Status(http.StatusBadRequest).SendString("An error occured!")
	}
	return ctx.Status(http.StatusOK).Send(userByte)
}
