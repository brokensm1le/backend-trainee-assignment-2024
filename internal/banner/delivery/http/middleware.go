package http

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *BannerHandler) userIdentity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		headers := c.GetReqHeaders()
		//"Authorization"
		header, ok := headers["Authorization"]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		headerParts := strings.Split(header[0], " ")
		if len(headerParts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
		}

		tokenData, err := h.tokenManager.Parse(headerParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
		}

		c.Locals("tokenData", tokenData)
		return c.Next()
	}
}

//func (s *BannerHandler) userIdentity(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//		header := r.Header.Get("Authorization")
//		if header == "" {
//			http.Error(rw, fmt.Sprintf("empty auth header"), http.StatusUnauthorized)
//			return
//		}
//
//		headerParts := strings.Split(header, " ")
//		if len(headerParts) != 2 {
//			http.Error(rw, fmt.Sprintf("invalid auth header"), http.StatusUnauthorized)
//			return
//		}
//
//		tokenData, err := s.authUC.ParseToken(headerParts[1])
//		if err != nil {
//			http.Error(rw, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		ctx := context.WithValue(r.Context(), "tokenData", tokenData)
//
//		h.ServeHTTP(rw, r.WithContext(ctx))
//	})
//}
