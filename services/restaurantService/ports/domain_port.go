package ports

import (
	"context"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

// DomainPort defines the interface for restaurant domain operations.
type DomainPort interface {
	// Restaurant-related methods
	GetAllRestaurantsDomain(ctx context.Context) ([]generated.Restaurant, error)
	GetRestaurantByIdDomain(ctx context.Context, restaurantId int32) (*generated.Restaurant, error)
	FilterRestaurantsByCategoryDomain(ctx context.Context, category string) ([]generated.Restaurant, error)

	// Menu-related methods
	GetMenuItemsByRestaurantIdDomain(ctx context.Context, restaurantId int32) ([]generated.Menuitem, error)
	GetMenuItemByRestaurantAndIdDomain(ctx context.Context, params generated.GetMenuItemByRestaurantAndIdParams) (*generated.Menuitem, error)

	// Category-related methods
	GetAllCategoriesDomain(ctx context.Context) ([]string, error)
}
