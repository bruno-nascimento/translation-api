package main_test

import (
	"context"
	"fmt"
	"io"
	"net"
	defaulthttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http"
	"github.com/bruno-nascimento/translation-api/internal/mocks"
	"github.com/bruno-nascimento/translation-api/internal/repository"
	"github.com/bruno-nascimento/translation-api/internal/service"

	"github.com/bruno-nascimento/translation-api/pkg/client"
)

func TestServer(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockQuerier(ctrl)

	t.Run("Should create a translation successfully using the openAPI client", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping testing in short mode")
		}
		ctx := context.Background()

		cfg, err := config.New(ctx)
		cfg.Server.Port = "3333"
		assert.NoError(t, err)

		translationAPI := http.NewTranslationAPI(service.NewTranslation(cfg, mockRepo))

		mockRepo.EXPECT().InsertWord(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		mockRepo.EXPECT().InsertTranslation(gomock.Any(), gomock.Any()).Return(nil)

		srv := http.NewServer(cfg, translationAPI)
		defer srv.Shutdown()

		go func() {
			if err = srv.Listen(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)); err != nil {
				assert.NoError(t, err)
			}
		}()

		var defaultClient = &defaulthttp.Client{}
		defaultClient.Timeout = time.Second * 3

		cli, err := client.NewClient(fmt.Sprintf("http://%s:%s", cfg.Server.Host, cfg.Server.Port), client.WithHTTPClient(defaultClient))
		assert.NoError(t, err)

		body := strings.NewReader(`{
			"from": {
				"word": "mesa",
					"language": "pt-br"
			},
			"to": {
				"word": "desk",
					"language": "en-US"
			}
		}`)

		response, err := cli.AddTranslationWithBody(ctx, fiber.MIMEApplicationJSON, body)

		time.Sleep(time.Second * 2)

		assert.NoError(t, err)
		assert.Equal(t, response.StatusCode, fiber.StatusCreated)
		respBody, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(respBody), "Created")
	})

	t.Run("Validation tests using fiber test support", func(t *testing.T) {
		tests := []struct {
			name            string
			route           string
			expectedCode    int
			prepareRepoMock func(mockRepo *mocks.MockQuerier) *mocks.MockQuerier
			expectedBody    string
		}{
			{
				name:         "No content - no translation or suggestions found",
				route:        "/v1/translation?word=nada&language=pt-br&target_language=en-US",
				expectedCode: 204,
				prepareRepoMock: func(mockRepo *mocks.MockQuerier) *mocks.MockQuerier {
					mockRepo.EXPECT().SelectTranslations(gomock.Any(), gomock.Any()).Return(nil, nil)
					mockRepo.EXPECT().SelectSimilarWords(gomock.Any(), gomock.Any()).Return(nil, nil)
					return mockRepo
				},
			},
			{
				name:         "bad request: invalid target language",
				route:        "/v1/translation?word=invalido&language=pt-br&target_language=invalid-target-language",
				expectedCode: 400,
				expectedBody: `{"code":400,"message":"language: tag is not well-formed"}`,
			},
			{
				name:         "bad request: invalid language",
				route:        "/v1/translation?word=invalido&language=ptbr&target_language=en-us",
				expectedCode: 400,
				expectedBody: `{"code":400,"message":"language: tag is not well-formed"}`,
			},
			{
				name:         "bad request: no word param provided",
				route:        "/v1/translation?language=pt-br&target_language=en-us",
				expectedCode: 500,                                              // TODO: should be 400 - there are some problems with the openapi generator
				expectedBody: `Query argument word is required, but not found`, // TODO should be a json error, not a plain text
			},
			{
				name:         "OK - translation found",
				route:        "/v1/translation?word=mesa&language=pt-br&target_language=en-US",
				expectedCode: 200,
				prepareRepoMock: func(mockRepo *mocks.MockQuerier) *mocks.MockQuerier {
					mockRepo.EXPECT().SelectTranslations(gomock.Any(), repository.SelectTranslationsParams{
						ToLang:         "en",
						ToLangRegion:   pgtype.Text{String: "US", Valid: true},
						FromWord:       "mesa",
						FromLang:       "pt",
						FromLangRegion: pgtype.Text{String: "BR", Valid: true},
					}).Return([]repository.Word{
						{ID: "ulid", Value: "desk", Lang: "en", LangRegion: pgtype.Text{String: "US", Valid: true}},
					}, nil)
					return mockRepo
				},
				expectedBody: `{"results":["desk"]}`,
			},
			{
				name:         "Not found - show suggestions",
				route:        "/v1/translation?word=mesa&language=pt-br&target_language=en-US",
				expectedCode: 404,
				prepareRepoMock: func(mockRepo *mocks.MockQuerier) *mocks.MockQuerier {
					mockRepo.EXPECT().SelectTranslations(gomock.Any(), repository.SelectTranslationsParams{
						ToLang:         "en",
						ToLangRegion:   pgtype.Text{String: "US", Valid: true},
						FromWord:       "mesa",
						FromLang:       "pt",
						FromLangRegion: pgtype.Text{String: "BR", Valid: true},
					}).Return(nil, nil)
					mockRepo.EXPECT().SelectSimilarWords(gomock.Any(), repository.SelectSimilarWordsParams{
						StrictWordSimilarity: "mesa",
						Lang:                 "pt",
						LangRegion:           pgtype.Text{String: "BR", Valid: true},
						Limit:                5,
					}).Return([]repository.SelectSimilarWordsRow{
						{Value: "missa", StrictWordSimilarity: 0.9},
						{Value: "mesada", StrictWordSimilarity: 0.8},
						{Value: "meio", StrictWordSimilarity: 0.5},
						{Value: "fora", StrictWordSimilarity: 0.0},
						{Value: "medo", StrictWordSimilarity: 0.0},
					}, nil)
					return mockRepo
				},
				expectedBody: `{"similar_words":["missa","mesada","meio"]}`,
			},
		}

		ctx := context.Background()

		cfg, err := config.New(ctx)
		assert.NoError(t, err)
		translationAPI := http.NewTranslationAPI(service.NewTranslation(cfg, mockRepo))
		srv := http.NewServer(cfg, translationAPI)
		defer srv.Shutdown()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if test.prepareRepoMock != nil {
					mockRepo = test.prepareRepoMock(mockRepo)
				}
				req := httptest.NewRequest("GET", test.route, nil)
				resp, _ := srv.Test(req, -1)
				assert.Equal(t, test.expectedCode, resp.StatusCode)
				respBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBody, string(respBody))
			})
		}
	})
}
