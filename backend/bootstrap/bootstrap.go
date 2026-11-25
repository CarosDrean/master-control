package bootstrap

import (
	"context"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"master-control/domain/action"
	"master-control/domain/actiontype"
	"master-control/domain/page"
	"master-control/infrastructure/handler"
	"master-control/infrastructure/middleware"
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

	cfg, err := loadConfig(getenv)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}

	db, err := newBbolt()
	if err != nil {
		log.Fatal().Err(err).Msg("initialize bbolt failed")
	}
	defer func() { _ = db.Close() }()

	api := gin.New()
	api.Use(gin.Recovery())
	// El scroll táctil llega varias veces por segundo.
	api.Use(middleware.Logger(path.Join("/", appName, "system/scroll")))
	api.Use(cors.New(corsConfig(cfg)))
	apiGroup := api.Group(appName)
	apiGroup.Use(middleware.Auth(cfg.Token, healthPath(appName)))

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
			ActionUseCase:     action.New(actionRepository.New(db), actionTypeUseCase, systemController, cfg.MaxTimeout),
			ActionTypeUseCase: actionTypeUseCase,
			PageUseCase:       pageUseCase,
			SystemUseCase:     systemController,
		},
	})

	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: api,
		// La escritura debe sobrevivir a la acción más larga permitida (y a la captura de hotkey, 15 s).
		WriteTimeout: max(cfg.MaxTimeout, 15*time.Second) + 10*time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if cfg.TokenGenerated {
		log.Warn().Msgf("MC_TOKEN vacío: token efímero para esta ejecución: %s (define MC_TOKEN en backend/.env para fijarlo)", cfg.Token)
	}
	if len(cfg.AllowedOrigins) == 0 {
		log.Warn().Msg("MC_ALLOWED_ORIGINS vacío: se aceptan todos los orígenes (el token sigue siendo obligatorio)")
	}

	log.Info().Msgf("Starting server at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}

// healthPath es la ruta completa del health check, la única pública.
func healthPath(appName string) string {
	return path.Join("/", appName, "health")
}

// corsConfig restringe los orígenes a MC_ALLOWED_ORIGINS; vacío = todos (con token obligatorio).
func corsConfig(cfg Config) cors.Config {
	c := cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}
	if len(cfg.AllowedOrigins) == 0 {
		c.AllowAllOrigins = true
	} else {
		c.AllowOrigins = cfg.AllowedOrigins
	}
	return c
}
