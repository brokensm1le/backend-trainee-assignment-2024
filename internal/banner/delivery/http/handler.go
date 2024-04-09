package http

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/pkg/tokenManager"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

type BannerHandler struct {
	bannerUC     banner.Usecase
	tokenManager tokenManager.TokenManager
}

func NewBannerHandler(bannerUC banner.Usecase, tokenManager tokenManager.TokenManager) *BannerHandler {
	return &BannerHandler{
		bannerUC:     bannerUC,
		tokenManager: tokenManager,
	}
}

func (h *BannerHandler) GetBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params banner.GetBannerParams
		)

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		log.Println("tokenData:", tokenData)

		if !validateParamsGetBanner(c) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
		}

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
		var (
			params banner.GetFilteredBannersParams = banner.GetFilteredBannersParams{
				TagID:     -1,
				FeatureID: -1,
				Limit:     5,
			}
		)

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		log.Println("tokenData:", tokenData)

		if err := c.QueryParser(&params); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if (params.TagID == -1 && params.FeatureID == -1) || params.Limit <= 0 || params.Offset < 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bad params"})
		}

		banners, err := h.bannerUC.GetFilteredBanners(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(banners)
	}
}

func (h *BannerHandler) CreateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params banner.CreateBannerParams
		)

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		fmt.Println(data, ok)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		if data.Role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not admin role"})
		}

		log.Println("tokenData:", tokenData)

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad body"})
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
		var (
			params banner.UpdateBannerParams
		)

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		fmt.Println(data, ok)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		if data.Role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not admin role"})
		}

		log.Println("tokenData:", tokenData)

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "bad id"})
		}

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad body"})
		}
		params.BannerID = id

		err = h.bannerUC.UpdateUser(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}

func (h *BannerHandler) DeleteBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		if data.Role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not admin role"})
		}

		log.Println("tokenData:", tokenData)

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		err = h.bannerUC.DeleteBanner(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}

// ------------------------------------------------------------------------------------------------------

func validateParamsGetBanner(c *fiber.Ctx) bool {
	if tagId := c.Query("tag_id"); tagId == "" {
		return false
	}
	if fId := c.Query("feature_id"); fId == "" {
		return false
	}
	return true
}
