package internal

import (
	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/core/auth"
	"github.com/HasanNugroho/starter-golang/internal/core/roles"
	"github.com/HasanNugroho/starter-golang/internal/core/users"
	"github.com/HasanNugroho/starter-golang/internal/shared/middleware"
	"github.com/HasanNugroho/starter-golang/internal/shared/modules"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ProductionEnv = "production"
	redisClient   *redis.Client
	mongodb       *mongo.Database
)

func AppsInit(router *echo.Echo) *app.Apps {
	// Initialize configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
	}

	// Initialize Logger
	logApps := config.InitLogger(appConfig)

	// Initialize Database if enabled
	if appConfig.DB.Enabled {
		mongodb, err = appConfig.DB.InitMongo()
		if err != nil {
			logApps.Fatal().Msg(err.Error())
			panic(1)
		}
	}

	// Initialize Redis if enabled
	if appConfig.Redis.Enabled {
		redisClient, err = appConfig.Redis.InitRedis()
		if err != nil {
			logApps.Fatal().Msg(err.Error())
			panic(1)
		}
		// defer config.ShutdownRedis(redisClient)
	}

	// Initialize Elastic if enabled
	if appConfig.Search.Enabled {
		err := appConfig.Search.SearchInit()
		if err != nil {
			logApps.Fatal().Msg(err.Error())
			panic(1)
		}
	}

	app := &app.Apps{
		Config: appConfig,
		Log:    logApps,
		DB:     mongodb,
		Redis:  redisClient,
		Bus:    modules.EventNew(),
		Router: router,
	}

	app.Router.Use(middleware.SetCORS(app.Config), middleware.SecurityMiddleware(app.Config))

	// Initialize Rate Limiter if enabled
	app.Router.Use(middleware.RateLimit(appConfig))

	// Initialize modules
	InitModules(app)

	return app
}

func InitModules(app *app.Apps) {
	app.RegisterFeature(users.NewUserModule(app))
	app.RegisterFeature(auth.NewAuthModule(app))
	app.RegisterFeature(roles.NewRoleModule(app))

	app.InitFeatures()
}
