package test

import (
	"backend-trainee-assignment-2024/internal/auth"
	http2 "backend-trainee-assignment-2024/internal/auth/delivery/http"
	mockAuth "backend-trainee-assignment-2024/internal/auth/mocks"
	"backend-trainee-assignment-2024/internal/auth/usecase"
	"backend-trainee-assignment-2024/internal/banner"
	http3 "backend-trainee-assignment-2024/internal/banner/delivery/http"
	mockBanner "backend-trainee-assignment-2024/internal/banner/mocks"
	usecase2 "backend-trainee-assignment-2024/internal/banner/usecase"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/hasher/SHA256"
	"backend-trainee-assignment-2024/pkg/tokenManager/jwtTokenManager"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	t.Run("testGetBanner", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockBannerRepo := mockBanner.NewMockRepository(c)
		content := "content!!!"
		mockBannerRepo.EXPECT().GetContentBannerAdmin(&banner.GetBannerParams{TagID: 2, FeatureID: 2, UseLastRevision: true, Role: 1}).Return(&content, nil).Times(1)

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BannerUC := usecase2.NewBannerUsecase(mockBannerRepo, nil)
		handlerBanner := http3.NewBannerHandler(BannerUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBanner.UserIdentity())
		appWithToken.Get("/user_banner", handlerBanner.GetBanner())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		json.NewEncoder(b).Encode(m)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		json.Unmarshal(bodyBytes, &tokens)

		// TEST get banner

		m = map[string]string{"tag_id": "2", "feature_id": "2", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/user_banner", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q := req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.Equal(t, content, string(bodyBytes))
	})
}
