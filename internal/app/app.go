package app

import (
	"fmt"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var GlobalApps *Apps

type Apps struct {
	Config   *config.Config
	Log      *zerolog.Logger
	Redis    *redis.Client
	DB       *gorm.DB
	Router   *gin.Engine
	Features []Feature
}

type Feature interface {
	Register(app *Apps) error
	Route(router *gin.RouterGroup, app *Apps)
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
