package domain

import (
	"context"
	"errors"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type RestaurantDomain struct {
	repo *generated.Queries
}

func (d *RestaurantDomain) GetAllRestaurantsDomain(ctx context.Context) ([]generated.Restaurant, error) {
	rows, err := d.repo.FetchAllRestaurants(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch restaurants")
	}

	var restaurants []generated.Restaurant
	for _, row := range rows {
		restaurants = append(restaurants, generated.Restaurant{
			ID:       row.ID,
			Name:     row.Name,
			Rating:   row.Rating,
			Category: row.Category,
			Address:  row.Address,
			ZipCode:  row.ZipCode,
		})
	}
	return restaurants, nil
}

func (d *RestaurantDomain) GetRestaurantByIdDomain(ctx context.Context, restaurantId int32) (*generated.Restaurant, error) {
	if restaurantId <= 0 {
		return nil, errors.New("invalid restaurant id")
	}

	row, err := d.repo.GetRestaurantById(ctx, restaurantId)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	restaurant := &generated.Restaurant{
		ID:       row.ID,
		Name:     row.Name,
		Rating:   row.Rating,
		Category: row.Category,
		Address:  row.Address,
		ZipCode:  row.ZipCode,
	}

	return restaurant, nil
}

func (d *RestaurantDomain) GetMenuItemsByRestaurantIdDomain(ctx context.Context, restaurantId int32) ([]generated.Menuitem, error) {
	if restaurantId <= 0 {
		return nil, errors.New("invalid restaurant id")
	}

	menuitems, err := d.repo.FetchMenuItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, errors.New("failed to fetch menuitems")
	}
	return menuitems, nil
}

func (d *RestaurantDomain) GetMenuItemByRestaurantAndIdDomain(ctx context.Context, params generated.GetMenuItemByRestaurantAndIdParams) (*generated.Menuitem, error) {
	row, err := d.repo.GetMenuItemByRestaurantAndId(ctx, params)
	if err != nil {
		return nil, errors.New("menuitem not found")
	}

	menuitem := &generated.Menuitem{
		ID:           row.ID,
		Restaurantid: row.Restaurantid,
		Name:         row.Name,
		Price:        row.Price,
		Description:  row.Description,
	}

	return menuitem, nil
}

func (d *RestaurantDomain) GetAllCategoriesDomain(ctx context.Context) ([]string, error) {
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

func (d *RestaurantDomain) FilterRestaurantsByCategoryDomain(ctx context.Context, category string) ([]generated.Restaurant, error) {
	if len(category) == 0 {
		return nil, errors.New("category cannot be empty")
	}

	rows, err := d.repo.FilterRestaurantsByCategory(ctx, &category)
	if err != nil {
		return nil, errors.New("failed to filter restaurants by category")
	}

	// Map rows to Restaurant objects
	var restaurants []generated.Restaurant
	for _, row := range rows {
		restaurants = append(restaurants, generated.Restaurant{
			ID:       row.ID,
			Name:     row.Name,
			Rating:   row.Rating,
			Category: row.Category,
			Address:  row.Address,
			ZipCode:  row.ZipCode,
		})
	}
	return restaurants, nil
}

// func (d *RestaurantDomain) PublishMenuItemSelection(ctx context.Context, queue string, event Event) error {
// 	err := d.broker.Publish(ctx, queue, event)
// 	if err != nil {
// 		return errors.New("failed to publish menu item selection")
// 	}
// 	return nil
// }
