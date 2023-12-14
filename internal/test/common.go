package test

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	FirstId  = uuid.MustParse("74c65c5a-44e7-44a2-8cdd-cd7c49cfcb42")
	SecondId = uuid.MustParse("a7670bf9-f4b0-4e5c-8edc-140812dbf719")
	ThirdId  = uuid.MustParse("45bdca7a-f8d8-42e5-9ad8-706a216647ab")
)

var Date = time.Date(2020, 1, 25, 1, 1, 1, 1, time.UTC)

func Must[T any](provider func() (T, error)) T {
	t, err := provider()
	if err != nil {
		log.Fatal().Err(err).Msg("failed run provider")
	}

	if IsNil(t) {
		log.Fatal().Msg("provider returned nil")
	}

	return t
}

func IsNil[T any](t T) bool {
	switch casted := any(t).(type) {
	case interface{}:
		return casted == nil
	default:
		return false
	}
}
