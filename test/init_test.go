package test

import (
	"testing"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	_, err := config.InitConfig()

	assert.Nil(t, err)
}
