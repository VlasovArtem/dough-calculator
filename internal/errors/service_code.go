package errors

import "github.com/google/uuid"

var (
	SourdoughRecipeNotFound = func(details string) error {
		return NewBadRequestError(10001, "sourdough not found", details)
	}
)
var (
	FlourByIdNotFound = func(id uuid.UUID) error {
		return FlourNotFound("flour with id " + id.String() + " not found")
	}
	FlourNotFound = func(details string) error {
		return NewBadRequestError(20001, "flour not found", details)
	}
)
