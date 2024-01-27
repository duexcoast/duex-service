package testgrp

import (
	"context"
	"net/http"

	"github.com/duexcoast/duex-service/foundation/web"
)

// Test is our example route
func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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
