package main

import (
	"testing"

	"github.com/yash27062004-lab/yash27062004-lab.github.io/internal/config"
)

func TestDefaultConfigIsValid(t *testing.T) {
	if err := config.Default().Validate(); err != nil {
		t.Fatalf("default configuration is invalid: %v", err)
	}
}
