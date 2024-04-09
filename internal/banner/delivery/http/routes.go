package http

import "github.com/gofiber/fiber/v2"

func MapRoutes(router fiber.Router, h *BannerHandler) {
	routerWithToken := router.Use("/", h.userIdentity())
	routerWithToken.Get("/user_banner", h.GetBanner())
	routerWithToken.Get("/banner", h.GetFilteredBanners())
	routerWithToken.Post("/banner", h.CreateBanner())
	routerWithToken.Patch("/banner/:id", h.UpdateBanner())
	routerWithToken.Delete("/banner/:id", h.DeleteBanner())
}
