package shared

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

func GetOrigin(ctx context.Context) string {
	return ctx.Value(entity.OriginKey).(string)
}

func GetOriginTypeKey(ctx context.Context) string {
	return ctx.Value(entity.OriginTypeKey).(string)
}

func GetOrganization(ctx context.Context) *entity.OrganizationData {
	return ctx.Value(entity.OrganizationKey).(*entity.OrganizationData)
}

func GetUserCredential(ctx context.Context) *entity.UserCredentialData {
	if value := ctx.Value(entity.UserCredentialDataKey); value != nil {
		if userCredential, ok := value.(*entity.UserCredentialData); ok {
			return userCredential
		}
	}
	return nil
}

func GetCompanyID(ctx context.Context) *string {
	if value := ctx.Value(entity.CompanyKey); value != nil {
		strValue := value.(string)
		return &strValue
	}
	return nil
}

func NewBackgroundContext(ctx context.Context) context.Context {
	newCtx := context.Background()

	contextKeys := []any{
		entity.OriginKey,
		entity.OriginTypeKey,
		entity.OrganizationKey,
		entity.UserContextKey,
		entity.UserCredentialDataKey,
		entity.CompanyKey,
	}

	for _, key := range contextKeys {
		if value := ctx.Value(key); value != nil {
			newCtx = context.WithValue(newCtx, key, value)
		}
	}

	return newCtx
}
