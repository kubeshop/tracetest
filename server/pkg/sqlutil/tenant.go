package sqlutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

func Tenant(ctx context.Context, query string, params ...any) (string, []any) {
	tenantID := TenantID(ctx)
	if tenantID == nil {
		return query, params
	}

	prefix := getQueryPrefix(query)

	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" %s tenant_id = $%d", prefix, paramNumber)

	return query + condition, append(params, *tenantID)
}

func TenantWithPrefix(ctx context.Context, query string, prefix string, params ...any) (string, []any) {
	tenantID := TenantID(ctx)
	if tenantID == nil {
		return query, params
	}

	queryPrefix := getQueryPrefix(query)
	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" %s %stenant_id = $%d", queryPrefix, prefix, paramNumber)

	return query + condition, append(params, *tenantID)
}

func TenantWithReplacedID(ctx context.Context, query string, id id.ID) (string, []any) {
	var params []any
	tenantID := TenantID(ctx)
	if tenantID == nil {
		return query, append(params, id)
	}

	prefix := getQueryPrefix(query)
	condition := fmt.Sprintf(" %s tenant_id = $%d", prefix, 2)

	return query + condition, append(params, *tenantID, *tenantID)
}

func TenantID(ctx context.Context) *string {
	tenantID := ctx.Value(middleware.TenantIDKey)

	if tenantID == "" || tenantID == nil {
		return nil
	}

	tenantIDString := tenantID.(string)
	return &tenantIDString
}

func getQueryPrefix(query string) string {
	prefix := ""
	if strings.Contains(strings.ToLower(query), "where") {
		prefix = "AND "
	} else {
		prefix = "WHERE "
	}

	return prefix
}
