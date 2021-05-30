package store_redis

import (
	"context"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
)

type InspectionRepository interface {
	/// ....
	GetListShowDataByEO(ctx context.Context, string string) (*model.Inspections, error)
	SetListShowDataByEO(ctx context.Context, n *model.Inspections) error
}
