package httpServer

import (
	http2 "backend-trainee-assignment-2024/internal/auth/delivery/http"
	authRepository "backend-trainee-assignment-2024/internal/auth/repository"
	authUsecase "backend-trainee-assignment-2024/internal/auth/usecase"
	"backend-trainee-assignment-2024/internal/banner/cache/lazyCache"
	"backend-trainee-assignment-2024/internal/banner/delivery/http"
	bannerRepository "backend-trainee-assignment-2024/internal/banner/repository"
	bannerUsecase "backend-trainee-assignment-2024/internal/banner/usecase"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/hasher/SHA256"
	"backend-trainee-assignment-2024/pkg/storagePostgres"
	"backend-trainee-assignment-2024/pkg/tokenManager/jwtTokenManager"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
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

	cache := lazyCache.NewCache(bannerRepo)

	authUC := authUsecase.NewAuthUsecase(authRepo, hasher, manager)
	bannerUC := bannerUsecase.NewBannerUsecase(bannerRepo, cache)

	authR := http2.NewAuthHandler(authUC)
	bannerR := http.NewBannerHandler(bannerUC, manager)

	http2.MapRoutes(app, authR)
	http.MapRoutes(app, bannerR)

	if s.cfg.Server.WithGenerator {
		err = storagePostgres.GenerateTable(s.cfg, bannerRepo)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	go func() {
		for {
			err = cache.LoadCache()
			if err != nil {
				log.Println("Error in loadCache:", err.Error())
			}
			log.Println("Cached Done!")
			time.Sleep(5 * time.Minute)
		}
	}()

	return nil
}
