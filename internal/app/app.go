package app

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/modules"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var GlobalApps *Apps

type Apps struct {
	Config   *config.Config
	Log      *zerolog.Logger
	Redis    *redis.Client
	DB       *mongo.Database
	Bus      *modules.EventBus
	Router   *echo.Echo
	Features []Feature
}

type Feature interface {
	Register(app *Apps) error
	Route(router *echo.Group, app *Apps)
}

func (a *Apps) RegisterFeature(f Feature) {
	a.Features = append(a.Features, f)
}

func (a *Apps) InitFeatures() error {
	for _, feature := range a.Features {
		if err := feature.Register(a); err != nil {
			return fmt.Errorf("failed to register feature: %w", err)
		}
	}

	api := a.Router.Group("/api")
	for _, feature := range a.Features {
		feature.Route(api, a)
	}

	return nil
}

func GetApps() *Apps {
	if GlobalApps == nil {
		panic("apps not initialized")
	}
	return GlobalApps
}
