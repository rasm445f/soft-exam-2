package domain

import (
	"context"
	"errors"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type RestaurantDomain struct {
	repo *generated.Queries
}

// NewRestaurantDomain initializes the domain layer
func NewRestaurantDomain(repo *generated.Queries) *RestaurantDomain {
	return &RestaurantDomain{repo: repo}
}

func (d *RestaurantDomain) FetchAllRestaurants(ctx context.Context) ([]generated.Restaurant, error) {
	rows, err := d.repo.FetchAllRestaurants(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch restaurants")
	}

	var restaurants []generated.Restaurant
	for _, row := range rows {
		restaurants = append(restaurants, generated.Restaurant{
			ID: row.ID,
			Name: row.Name,
			Address: row.Address,
			Rating: row.Rating,
			Category: row.Category,
		})
	}
	return restaurants, nil
}

func (d *RestaurantDomain) GetRestaurantById(ctx context.Context, restaurantId int32) (*generated.Restaurant, error) {
	if restaurantId <= 0 {
		return nil, errors.New("invalid restaurant id")
	}

	row, err := d.repo.GetRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	restaurant := &generated.Restaurant{
		ID: row.ID,
		Name: row.Name,
		Address: row.Address,
		Rating: row.Rating,
		Category: row.Category,
	}

	return restaurant, nil
}

func (d *RestaurantDomain) FetchMenuItemsByRestaurantId(ctx context.Context, restaurantId int32) ([]generated.FetchMenuItemsByRestaurantIdRow, error) {
	if restaurantId <= 0 {
		return nil, errors.New("invalid restaurant id")
	}

	menuitems, err := d.repo.FetchMenuItemsByRestaurantId(ctx, restaurantId) 
	if err != nil {
		return nil, errors.New("failed to fetch menuitems")
	}
	return menuitems, nil
}

func (d *RestaurantDomain) GetMenuItemByRestaurantAndId(ctx context.Context, params generated.GetMenuItemByRestaurantAndIdParams) (*generated.Menuitem, error) {
	row, err := d.repo.GetMenuItemByRestaurantAndId(ctx, params)
	if err != nil {
		return nil, errors.New("menuitem not found")
	}

	menuitem := &generated.Menuitem{
		ID: row.ID,
		Name: row.Name,
		Description: row.Description,
		Price: row.Price,
		Restaurantid: row.Restaurantid,
	}

	return menuitem, nil
} 

func (d* RestaurantDomain) FetchAllCategories(ctx context.Context) ([]string, error) {
	rows, err := d.repo.FetchAllCategories(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch categories")
	}

	var categories []string
	for _, row := range rows {
		if row != nil {
			categories = append(categories, *row)
		}
	}

	return categories, nil
}

func (d* RestaurantDomain) FilterRestaurantsByCategory(ctx context.Context, category string) ([]generated.Restaurant, error) {
	if len(category) == 0 {
		return nil, errors.New("category cannot be empty")
	}

	restaurants, err := d.repo.FilterRestaurantsByCategory(ctx, &category)
	if err != nil {
		return nil, errors.New("failed to filter restaurants by category")
	}
	return restaurants, nil
}

