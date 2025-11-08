package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVideoAspectRatio(t *testing.T) {
	tests := []struct {
		name       string
		w          int
		h          int
		expectedAR string
	}{
		{
			name:       "landscape",
			w:          1280,
			h:          720,
			expectedAR: "16:9",
		}, {
			name:       "portrait",
			w:          608,
			h:          1080,
			expectedAR: "9:16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aspectRatio := calculateAspectRatio(tt.w, tt.h)
			assert.Equal(t, tt.expectedAR, aspectRatio)
		})
	}
}
