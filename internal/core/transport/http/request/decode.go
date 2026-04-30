package core_http_request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if dest == nil {
		return fmt.Errorf("decode json: dest is nil")
	}
	if reflect.ValueOf(dest).Kind() != reflect.Pointer {
		return fmt.Errorf("decode json: dest must be a pointer")
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	// Prevent silent acceptance of multiple JSON values (e.g. `{} {}`)
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode json: unexpected trailing data")
		}
		return fmt.Errorf("decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("request validation: %w", err)
	}

	return nil
}
