package http

import "github.com/gofiber/fiber/v2"

func MapRoutes(router fiber.Router, h *BannerHandler) {
	router.Get("/user_banner", h.GetBanner())
	router.Get("/banner", h.GetFilteredBanners())
	router.Post("/banner", h.CreateBanner())
	router.Patch("/banner/:id", h.UpdateBanner())
	router.Delete("/banner/:id", h.DeleteBanner())
}
