package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type key string

var (
	TenantIDKey key = "tenantID"
)

const HeaderTenantID = "X-Tracetest-TenantID"

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenantID := getTenantIDFromRequest(r)

		// if tenant id exists and is invalid we return a 400 error
		if tenantID != "" && !isValidUUID(tenantID) {
			err := fmt.Errorf("invalid tenant id: %s", tenantID)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		ctx = context.WithValue(ctx, TenantIDKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TenantIDFromContext(ctx context.Context) string {
	tenantID := ctx.Value(TenantIDKey)

	if tenantID == nil {
		return ""
	}

	return tenantID.(string)
}

func getTenantIDFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderTenantID)
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
