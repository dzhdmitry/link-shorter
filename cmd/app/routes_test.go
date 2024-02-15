package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoutes(t *testing.T) {
	app := Application{}

	assert.NotNil(t, app.routes())
}
