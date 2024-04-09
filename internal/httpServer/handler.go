package httpServer

import (
	http2 "backend-trainee-assignment-2024/internal/auth/delivery/http"
	authRepository "backend-trainee-assignment-2024/internal/auth/repository"
	authUsecase "backend-trainee-assignment-2024/internal/auth/usecase"
	"backend-trainee-assignment-2024/internal/banner/delivery/http"
	bannerRepository "backend-trainee-assignment-2024/internal/banner/repository"
	bannerUsecase "backend-trainee-assignment-2024/internal/banner/usecase"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/hasher/SHA256"
	"backend-trainee-assignment-2024/pkg/storagePostgres"
	"backend-trainee-assignment-2024/pkg/tokenManager/jwtTokenManager"
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

	hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
	manager, err := jwtTokenManager.NewManger(cconstant.SignedKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	authRepo := authRepository.NewPostgresRepository(db)
	bannerRepo := bannerRepository.NewPostgresRepository(db)

	authUC := authUsecase.NewAuthUsecase(authRepo, hasher, manager)
	bannerUC := bannerUsecase.NewBannerUsecase(bannerRepo)

	authR := http2.NewAuthHandler(authUC)
	bannerR := http.NewBannerHandler(bannerUC, manager)

	http2.MapRoutes(app, authR)
	http.MapRoutes(app, bannerR)

	return nil
}
