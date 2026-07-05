package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeResponse(t *testing.T) {
	resp := parseLlamaGuardResp("safe")

	assert.True(t, resp.Safe, "prompt should be safe")
}

func TestUnsafeResponse(t *testing.T) {
	resp := parseLlamaGuardResp("unsafe\nS2")

	assert.False(t, resp.Safe, "prompt should be not safe")
	assert.Equal(t, "Non-Violent Crimes", resp.Reason)
}
