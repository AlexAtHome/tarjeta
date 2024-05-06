package auth

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginDetails struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
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
	cookie.HTTPOnly = true
	cookie.SameSite = "lax"

	return cookie
}

// POST
func Login(ctx *fiber.Ctx) error {
	p := new(LoginDetails)
	if err := ctx.BodyParser(p); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(string("An error occured!"))
	}

	// TODO: Add database
	if p.Username != "alex" && p.Password != "qwerty123$!" {
		return ctx.Status(http.StatusUnauthorized).SendString(string("Wrong credentials"))
	}

	session := constructCookie("ses", "a51da2e0-31a7-40e5-a102-490c4b8647eb", true)
	user := constructCookie("user", p.Username, true)
	user.HTTPOnly = false
	ctx.Cookie(session)
	ctx.Cookie(user)

	return ctx.Redirect("/profile.html", http.StatusMovedPermanently)
}

// POST
func Logout(ctx *fiber.Ctx) error {
	session := constructCookie("ses", "", true)
	session.SessionOnly = false
	session.Expires = time.Now().Add(-(time.Hour * 2))
	user := constructCookie("user", "", false)
	user.SessionOnly = false
	user.Expires = time.Now().Add(-(time.Hour * 2))
	ctx.Cookie(session)
	ctx.Cookie(user)
	return ctx.Redirect("/login.html", http.StatusMovedPermanently)
}
