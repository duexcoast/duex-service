package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/duexcoast/duex-service/business/web/v1"
	"github.com/duexcoast/duex-service/foundation/web"
)

// Test is our example route
func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return v1.NewRequestError(errors.New("TRUSTED ERROR"), http.StatusBadRequest)
	}

	// All our handlers should be doing:
	// - Validate the data
	// - Call into the business layer
	// - Return errors
	// - Handle OK response
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	web.Respond(ctx, w, status, http.StatusOK)
	return nil
}
