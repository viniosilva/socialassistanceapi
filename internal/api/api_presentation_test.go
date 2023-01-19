package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ApiPresentation_BuildPreviousURL(t *testing.T) {
	cases := map[string]struct {
		inputHost   string
		inputLimit  int
		inputOffset int
		expectedUrl string
	}{
		"should return url with offset when offset is greater than zero": {
			inputHost:   "http://localhost:8080",
			inputLimit:  10,
			inputOffset: 20,
			expectedUrl: "http://localhost:8080?limit=10&offset=10",
		},
		"should return url without offset when offset is zero": {
			inputHost:   "http://localhost:8080",
			inputLimit:  10,
			inputOffset: 10,
			expectedUrl: "http://localhost:8080?limit=10",
		},
		"should return empty url when offset is zero": {
			inputHost:   "http://localhost:8080",
			inputLimit:  10,
			inputOffset: 0,
			expectedUrl: "",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			url := BuildPreviousURL(cs.inputHost, cs.inputLimit, cs.inputOffset)

			// then
			assert.Equal(t, cs.expectedUrl, url)
		})
	}
}

func Test_ApiPresentation_BuildNextURL(t *testing.T) {
	cases := map[string]struct {
		inputHost   string
		inputLimit  int
		inputOffset int
		inputTotal  int
		expectedUrl string
	}{
		"should return url when total is greater than limit": {
			inputHost:   "http://localhost:8080",
			inputLimit:  10,
			inputOffset: 0,
			inputTotal:  100,
			expectedUrl: "http://localhost:8080?limit=10&offset=10",
		},
		"should return empty url when total is equal to offset + limit": {
			inputHost:   "http://localhost:8080",
			inputLimit:  10,
			inputOffset: 0,
			inputTotal:  10,
			expectedUrl: "",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			url := BuildNextURL(cs.inputHost, cs.inputLimit, cs.inputOffset, cs.inputTotal)

			// then
			assert.Equal(t, cs.expectedUrl, url)
		})
	}
}
