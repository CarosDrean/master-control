package bootstrap

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"master-control/domain/action"
	"master-control/domain/actiontype"
	"master-control/domain/page"
	"master-control/infrastructure/handler"
	actionRepository "master-control/infrastructure/repository/action"
	actiontypeRepository "master-control/infrastructure/repository/actiontype"
	pageRepository "master-control/infrastructure/repository/page"
	"master-control/infrastructure/system"
	"master-control/model"
)

func Run() {
	// we ignore the error because we want to load the .env file only in development
	_ = godotenv.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configLogger()

	gin.SetMode(os.Getenv("GIN_MODE"))

	appName := os.Getenv("APP_NAME")

	db, err := newBbolt()
	if err != nil {
		log.Fatal().Err(err).Msg("initialize bbolt failed")
	}
	defer func() { _ = db.Close() }()

	api := gin.New()
	api.Use(gin.Recovery())
	api.Use(gin.Logger())
	api.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))
	apiGroup := api.Group(appName)

	actionTypeRepo := actiontypeRepository.New(db)
	if err := actionTypeRepo.Seed(ctx); err != nil {
		log.Fatal().Err(err).Msg("seed action types failed")
	}

	actionTypeUseCase := actiontype.New(actionTypeRepo)

	pageRepo := pageRepository.New(db)
	if err := pageRepo.Seed(ctx); err != nil {
		log.Fatal().Err(err).Msg("seed pages failed")
	}
	pageUseCase := page.New(pageRepo)

	systemController := system.New()

	handler.InitRoutes(model.RouterSpecification{
		Api: apiGroup,
		Dependencies: model.Dependencies{
			ActionUseCase:     action.New(actionRepository.New(db), actionTypeUseCase, systemController),
			ActionTypeUseCase: actionTypeUseCase,
			PageUseCase:       pageUseCase,
			SystemUseCase:     systemController,
		},
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      api,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msgf("Starting server at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
