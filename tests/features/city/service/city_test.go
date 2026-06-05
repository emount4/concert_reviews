package city_service_test

import (
	"context"
	"errors"
	"testing"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	city_service "github.com/emount4/concert_reviews/internal/features/city/service"
)

type fakeCityRepository struct {
	createCityFn   func(ctx context.Context, city core_models.City) (core_models.City, error)
	getCitiesFn    func(ctx context.Context, limit, offset *int) ([]core_models.City, error)
	getCityByIDFn  func(ctx context.Context, id int) (core_models.City, error)
	deleteCityFn   func(ctx context.Context, id int) error
	hasVenuesFn    func(ctx context.Context, id int) (bool, error)
	updateCityFn   func(ctx context.Context, city core_models.City) (core_models.City, error)

	createCalls  int
	getCalls     int
	getByIDCalls int
	deleteCalls  int
	hasCalls     int
	updateCalls  int

	lastCity   core_models.City
	lastID     int
	lastLimit  *int
	lastOffset *int
}

func (f *fakeCityRepository) CreateCity(ctx context.Context, city core_models.City) (core_models.City, error) {
	f.createCalls++
	f.lastCity = city
	if f.createCityFn != nil {
		return f.createCityFn(ctx, city)
	}
	city.CityID = 1
	return city, nil
}

func (f *fakeCityRepository) GetCities(ctx context.Context, limit, offset *int) ([]core_models.City, error) {
	f.getCalls++
	f.lastLimit = limit
	f.lastOffset = offset
	if f.getCitiesFn != nil {
		return f.getCitiesFn(ctx, limit, offset)
	}
	return []core_models.City{}, nil
}

func (f *fakeCityRepository) GetCityByID(ctx context.Context, id int) (core_models.City, error) {
	f.getByIDCalls++
	f.lastID = id
	if f.getCityByIDFn != nil {
		return f.getCityByIDFn(ctx, id)
	}
	return core_models.City{CityID: id, Name: "City", Slug: "city", Timezone: "UTC+3"}, nil
}

func (f *fakeCityRepository) DeleteCity(ctx context.Context, id int) error {
	f.deleteCalls++
	f.lastID = id
	if f.deleteCityFn != nil {
		return f.deleteCityFn(ctx, id)
	}
	return nil
}

func (f *fakeCityRepository) HasVenues(ctx context.Context, id int) (bool, error) {
	f.hasCalls++
	f.lastID = id
	if f.hasVenuesFn != nil {
		return f.hasVenuesFn(ctx, id)
	}
	return false, nil
}

func (f *fakeCityRepository) UpdateCity(ctx context.Context, city core_models.City) (core_models.City, error) {
	f.updateCalls++
	f.lastCity = city
	if f.updateCityFn != nil {
		return f.updateCityFn(ctx, city)
	}
	return city, nil
}

func TestCityServiceCreateGeneratesSlug(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)

	city := core_models.City{Name: "Санкт-Петербург", Timezone: "UTC+3"}
	created, err := service.Create(context.Background(), city)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected CreateCity once, got %d", repo.createCalls)
	}
	if repo.lastCity.Slug == "" {
		t.Fatal("expected generated slug to be set")
	}
	if repo.lastCity.Slug != "sankt-peterburg" {
		t.Fatalf("expected slug 'sankt-peterburg', got %q", repo.lastCity.Slug)
	}
	if created.Slug != "sankt-peterburg" {
		t.Fatalf("expected created slug 'sankt-peterburg', got %q", created.Slug)
	}
}

func TestCityServiceCreateKeepsProvidedSlug(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)

	city := core_models.City{Name: "Moscow", Slug: "custom-slug", Timezone: "UTC+3"}
	_, err := service.Create(context.Background(), city)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}
	if repo.lastCity.Slug != "custom-slug" {
		t.Fatalf("expected slug to remain 'custom-slug', got %q", repo.lastCity.Slug)
	}
}

func TestCityServiceCreateRejectsInvalidCity(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)

	_, err := service.Create(context.Background(), core_models.City{Name: "City"})
	if err == nil {
		t.Fatal("Create() expected error for invalid city, got nil")
	}
}

func TestCityServiceGetCitiesRejectsNegativeLimit(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)
	limit := -1

	_, err := service.GetCities(context.Background(), &limit, nil)
	if err == nil {
		t.Fatal("GetCities() expected error for negative limit, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("GetCities() expected invalid argument, got %v", err)
	}
}

func TestCityServiceGetCitiesRejectsNegativeOffset(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)
	offset := -5

	_, err := service.GetCities(context.Background(), nil, &offset)
	if err == nil {
		t.Fatal("GetCities() expected error for negative offset, got nil")
	}
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("GetCities() expected invalid argument, got %v", err)
	}
}

func TestCityServiceDeleteRejectsWhenHasVenues(t *testing.T) {
	repo := &fakeCityRepository{
		hasVenuesFn: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
	}
	service := city_service.NewCityService(repo)

	err := service.Delete(context.Background(), 5)
	if err == nil {
		t.Fatal("Delete() expected conflict error, got nil")
	}
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("Delete() expected conflict, got %v", err)
	}
	if repo.deleteCalls != 0 {
		t.Fatalf("expected no delete call when conflict, got %d", repo.deleteCalls)
	}
}

func TestCityServiceUpdateAppliesChanges(t *testing.T) {
	repo := &fakeCityRepository{}
	service := city_service.NewCityService(repo)

	name := "New City"
	slug := "new-city"
	tz := "UTC+3"

	updated, err := service.Update(context.Background(), 1, &name, &slug, &tz)
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}
	if repo.getByIDCalls != 1 || repo.updateCalls != 1 {
		t.Fatalf("expected GetCityByID and UpdateCity, got get=%d update=%d", repo.getByIDCalls, repo.updateCalls)
	}
	if repo.lastCity.Name != "New City" || repo.lastCity.Slug != "new-city" || repo.lastCity.Timezone != "UTC+3" {
		t.Fatalf("Update() unexpected city state: %+v", repo.lastCity)
	}
	if updated.Name != "New City" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
}
