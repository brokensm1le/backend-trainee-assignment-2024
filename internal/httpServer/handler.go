package httpServer

import (
	"backend-trainee-assignment-2024/internal/banner/delivery/http"
	"backend-trainee-assignment-2024/internal/banner/repository"
	"backend-trainee-assignment-2024/internal/banner/usecase"
	"backend-trainee-assignment-2024/pkg/storagePostgres"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (s *Server) MapHandlers(app *fiber.App) error {

	db, err := storagePostgres.InitPsqlDB(s.cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = storagePostgres.CreateTable(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bannerRepo := repository.NewPostgresRepository(db)

	bannerUC := usecase.NewBannerUsecase(bannerRepo)

	bannerR := http.NewBannerHandler(bannerUC)

	http.MapRoutes(app, bannerR)

	return nil
}
