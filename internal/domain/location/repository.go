package location

import "context"

type Repository interface {
	Create(ctx context.Context, loc *Location) error
	GetByID(ctx context.Context, id string) (*Location, error)
	List(ctx context.Context, includeInactive bool) ([]*Location, error)
	Update(ctx context.Context, loc *Location) error
}
