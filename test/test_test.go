package test

import (
	"backend-trainee-assignment-2024/internal/auth"
	http2 "backend-trainee-assignment-2024/internal/auth/delivery/http"
	mockAuth "backend-trainee-assignment-2024/internal/auth/mocks"
	"backend-trainee-assignment-2024/internal/auth/usecase"
	"backend-trainee-assignment-2024/internal/banner"
	http3 "backend-trainee-assignment-2024/internal/banner/delivery/http"
	"backend-trainee-assignment-2024/internal/banner/fake_repository"
	mockBanner "backend-trainee-assignment-2024/internal/banner/mocks"
	usecase2 "backend-trainee-assignment-2024/internal/banner/usecase"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/hasher/SHA256"
	"backend-trainee-assignment-2024/pkg/tokenManager/jwtTokenManager"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type CreateBannerID struct {
	BannerID int64 `json:"banner_id"`
}

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
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

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
		require.NoError(t, err)
		require.Equal(t, content, string(bodyBytes))
	})

	t.Run("testNoAdmin", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 0}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBannerRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BannerUC := usecase2.NewBannerUsecase(FakeBannerRepo, nil)
		handlerBanner := http3.NewBannerHandler(BannerUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBanner.UserIdentity())
		appWithToken.Post("/banner", handlerBanner.CreateBanner())
		appWithToken.Get("/user_banner", handlerBanner.GetBanner())
		appWithToken.Get("/banner", handlerBanner.GetFilteredBanners())
		appWithToken.Delete("/banner/:id", handlerBanner.DeleteBanner())
		appWithToken.Patch("/banner/:id", handlerBanner.UpdateBanner())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create banner
		r := banner.CreateBannerParams{TagIDs: []int64{2, 3}, FeatureID: 2, Content: "content!!!", IsActive: true, UseLastRevision: true}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/banner", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, "403 Forbidden", resp.Status)

		// TEST delete banner
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/banner/%d", 1), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, "403 Forbidden", resp.Status)

		// TEST update banner
		req, _ = http.NewRequest("PATCH", fmt.Sprintf("/banner/%d", 1), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, "403 Forbidden", resp.Status)

		// TEST get banner content
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
		require.NoError(t, err)
		require.Equal(t, "{\"error\":\"sql: no rows in result set\"}", string(bodyBytes))
	})

	t.Run("testGetBannerWithFake", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBannerRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BannerUC := usecase2.NewBannerUsecase(FakeBannerRepo, nil)
		handlerBanner := http3.NewBannerHandler(BannerUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBanner.UserIdentity())
		appWithToken.Post("/banner", handlerBanner.CreateBanner())
		appWithToken.Get("/user_banner", handlerBanner.GetBanner())
		appWithToken.Get("/banner", handlerBanner.GetFilteredBanners())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create banner
		r := banner.CreateBannerParams{TagIDs: []int64{2, 3}, FeatureID: 2, Content: "content!!!", IsActive: true, UseLastRevision: true}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/banner", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		bannerId := CreateBannerID{}
		err = json.Unmarshal(bodyBytes, &bannerId)
		require.NoError(t, err)
		require.NotEqual(t, 0, bannerId.BannerID)
		//fmt.Println(bannerId.BannerID)

		// TEST get banner content

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
		require.NoError(t, err)
		require.Equal(t, r.Content, string(bodyBytes))

		// TEST get all info banner

		m = map[string]string{"tag_id": "2", "feature_id": "2", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/banner", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q = req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res []banner.GetFilteredBannersResponse
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, r.Content, res[0].Content)
		require.Equal(t, r.FeatureID, res[0].FeatureID)
		require.Equal(t, r.IsActive, res[0].IsActive)
	})

	t.Run("testGetBanners", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBannerRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BannerUC := usecase2.NewBannerUsecase(FakeBannerRepo, nil)
		handlerBanner := http3.NewBannerHandler(BannerUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBanner.UserIdentity())
		appWithToken.Post("/banner", handlerBanner.CreateBanner())
		appWithToken.Get("/user_banner", handlerBanner.GetBanner())
		appWithToken.Get("/banner", handlerBanner.GetFilteredBanners())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		sliceBanners := []struct {
			ban         banner.CreateBannerParams
			expectedErr error
		}{
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1, 2, 3},
					FeatureID:       1,
					Content:         "Content(1)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1, 2, 3},
					FeatureID:       2,
					Content:         "Content(2)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{4},
					FeatureID:       1,
					Content:         "Content(3)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1},
					FeatureID:       3,
					Content:         "Content(4)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
		}

		for _, ban := range sliceBanners {
			b = new(bytes.Buffer)
			err := json.NewEncoder(b).Encode(ban.ban)
			require.NoError(t, err)
			req, _ = http.NewRequest("POST", "/banner", b)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err = app.Test(req, -1)
			require.NoError(t, err)
			bodyBytes, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			bannerId := CreateBannerID{}
			err = json.Unmarshal(bodyBytes, &bannerId)
			require.NoError(t, err)
			require.NotEqual(t, 0, bannerId.BannerID)
			fmt.Println("ID:", bannerId.BannerID)
		}

		// TEST only tag_id
		m = map[string]string{"tag_id": "2", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/banner", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q := req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res []banner.GetFilteredBannersResponse
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, 2, len(res))

		// TEST only feature_id
		m = map[string]string{"feature_id": "1", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/banner", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q = req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		res = []banner.GetFilteredBannersResponse{}
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, 2, len(res))
	})

	t.Run("testDeleteBanner", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBannerRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BannerUC := usecase2.NewBannerUsecase(FakeBannerRepo, nil)
		handlerBanner := http3.NewBannerHandler(BannerUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBanner.UserIdentity())
		appWithToken.Post("/banner", handlerBanner.CreateBanner())
		appWithToken.Get("/user_banner", handlerBanner.GetBanner())
		appWithToken.Get("/banner", handlerBanner.GetFilteredBanners())
		appWithToken.Delete("/banner/:id", handlerBanner.DeleteBanner())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		sliceBanners := []struct {
			ban         banner.CreateBannerParams
			expectedErr error
		}{
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1, 2, 3},
					FeatureID:       1,
					Content:         "Content(1)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1, 2, 3},
					FeatureID:       2,
					Content:         "Content(2)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{4},
					FeatureID:       1,
					Content:         "Content(3)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
			{
				ban: banner.CreateBannerParams{
					TagIDs:          []int64{1},
					FeatureID:       3,
					Content:         "Content(4)!",
					IsActive:        true,
					UseLastRevision: true,
				},
				expectedErr: nil,
			},
		}

		for _, ban := range sliceBanners {
			b = new(bytes.Buffer)
			err := json.NewEncoder(b).Encode(ban.ban)
			require.NoError(t, err)
			req, _ = http.NewRequest("POST", "/banner", b)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
			req.Header.Set("Content-Type", "application/json")
			resp, err = app.Test(req, -1)
			require.NoError(t, err)
			bodyBytes, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			bannerId := CreateBannerID{}
			err = json.Unmarshal(bodyBytes, &bannerId)
			require.NoError(t, err)
			require.NotEqual(t, 0, bannerId.BannerID)
			fmt.Println("ID:", bannerId.BannerID)
		}

		// TEST tag_id
		m = map[string]string{"tag_id": "2", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/banner", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q := req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res []banner.GetFilteredBannersResponse
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, 2, len(res))

		// TEST delete
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/banner/%d", res[0].BannerID), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		_, err = app.Test(req, -1)
		require.NoError(t, err)

		// TEST delete tag_id
		m = map[string]string{"tag_id": "2", "use_last_revision": "true"}
		req, _ = http.NewRequest("GET", "/banner", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		q = req.URL.Query()
		for key, val := range m {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		res = []banner.GetFilteredBannersResponse{}
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, 1, len(res))
	})

}
