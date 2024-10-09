package backup

import "context"

type Creator interface {
	Create(ctx context.Context) error
}
