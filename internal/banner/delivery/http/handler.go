package http

import (
	"backend-trainee-assignment-2024/internal/banner"
	"github.com/gofiber/fiber/v2"
)

type BannerHandler struct {
	bannerUC banner.Usecase
}

func NewBannerHandler(bannerUC banner.Usecase) *BannerHandler {
	return &BannerHandler{
		bannerUC: bannerUC,
	}
}

func (h *BannerHandler) GetBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params banner.GetBannerParams
		)

		if err := c.QueryParser(&params); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		content, err := h.bannerUC.GetBanner(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if _, err = c.Status(fiber.StatusOK).Response().BodyWriter().Write([]byte(*content)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return nil
	}
}

func (h *BannerHandler) GetFilteredBanners() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) CreateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params banner.CreateBannerParams
		)

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Bad request"})
		}

		bannerId, err := h.bannerUC.CreateBanner(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"banner_id": bannerId})
	}
}

func (h *BannerHandler) UpdateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) DeleteBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
