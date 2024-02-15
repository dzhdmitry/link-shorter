package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNow(t *testing.T) {
	c := Clock{}

	assert.NotNil(t, c.Now())
}
