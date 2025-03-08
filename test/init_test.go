package test

import (
	"testing"

	config "github.com/HasanNugroho/starter-golang/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	_, err := config.InitConfig()

	assert.Nil(t, err)
}
