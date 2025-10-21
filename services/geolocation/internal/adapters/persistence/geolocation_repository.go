package persistence

import (
	"strings"

	"gorm.io/gorm"

	"theraclosure/geolocation-service/internal/core/domain"
	"theraclosure/geolocation-service/internal/core/ports"
)

type GeolocationRepository struct {
	db *gorm.DB
}

func NewGeolocationRepository(database *Database) ports.GeolocationRepository {
	return &GeolocationRepository{
		db: database.GetDB(),
	}
}

// Country operations
func (r *GeolocationRepository) CreateCountry(country *domain.Country) error {
	return r.db.Create(country).Error
}

func (r *GeolocationRepository) GetCountryByID(id string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.Where("id = ?", id).First(&country).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}

func (r *GeolocationRepository) GetCountryByCode(code string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.Where("code = ? OR code2 = ?", strings.ToUpper(code), strings.ToUpper(code)).First(&country).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}

func (r *GeolocationRepository) GetAllCountries(limit, offset int) ([]*domain.Country, int64, error) {
	var countries []*domain.Country
	var total int64

	// Get total count
	if err := r.db.Model(&domain.Country{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get countries with pagination
	err := r.db.Limit(limit).Offset(offset).Order("name ASC").Find(&countries).Error
	if err != nil {
		return nil, 0, err
	}

	return countries, total, nil
}

func (r *GeolocationRepository) UpdateCountry(country *domain.Country) error {
	return r.db.Save(country).Error
}

func (r *GeolocationRepository) DeleteCountry(id string) error {
	return r.db.Delete(&domain.Country{}, "id = ?", id).Error
}

func (r *GeolocationRepository) SearchCountries(query string, limit, offset int) ([]*domain.Country, int64, error) {
	var countries []*domain.Country
	var total int64

	dbQuery := r.db.Model(&domain.Country{})

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(code2) LIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get countries with pagination
	err := dbQuery.Limit(limit).Offset(offset).Order("name ASC").Find(&countries).Error
	if err != nil {
		return nil, 0, err
	}

	return countries, total, nil
}

// State operations
func (r *GeolocationRepository) CreateState(state *domain.State) error {
	return r.db.Create(state).Error
}

func (r *GeolocationRepository) GetStateByID(id string) (*domain.State, error) {
	var state domain.State
	err := r.db.Preload("Country").Where("id = ?", id).First(&state).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &state, nil
}

func (r *GeolocationRepository) GetStatesByCountryID(countryID string, limit, offset int) ([]*domain.State, int64, error) {
	var states []*domain.State
	var total int64

	// Get total count
	if err := r.db.Model(&domain.State{}).Where("country_id = ?", countryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get states with pagination
	err := r.db.Preload("Country").Where("country_id = ?", countryID).
		Limit(limit).Offset(offset).Order("name ASC").Find(&states).Error
	if err != nil {
		return nil, 0, err
	}

	return states, total, nil
}

func (r *GeolocationRepository) UpdateState(state *domain.State) error {
	return r.db.Save(state).Error
}

func (r *GeolocationRepository) DeleteState(id string) error {
	return r.db.Delete(&domain.State{}, "id = ?", id).Error
}

func (r *GeolocationRepository) SearchStates(query string, countryID string, limit, offset int) ([]*domain.State, int64, error) {
	var states []*domain.State
	var total int64

	dbQuery := r.db.Model(&domain.State{}).Preload("Country")

	if countryID != "" {
		dbQuery = dbQuery.Where("country_id = ?", countryID)
	}

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", 
			searchPattern, searchPattern)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get states with pagination
	err := dbQuery.Limit(limit).Offset(offset).Order("name ASC").Find(&states).Error
	if err != nil {
		return nil, 0, err
	}

	return states, total, nil
}

// City operations
func (r *GeolocationRepository) CreateCity(city *domain.City) error {
	return r.db.Create(city).Error
}

func (r *GeolocationRepository) GetCityByID(id string) (*domain.City, error) {
	var city domain.City
	err := r.db.Preload("State").Preload("State.Country").Where("id = ?", id).First(&city).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &city, nil
}

func (r *GeolocationRepository) GetCitiesByStateID(stateID string, limit, offset int) ([]*domain.City, int64, error) {
	var cities []*domain.City
	var total int64

	// Get total count
	if err := r.db.Model(&domain.City{}).Where("state_id = ?", stateID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get cities with pagination
	err := r.db.Preload("State").Preload("State.Country").Where("state_id = ?", stateID).
		Limit(limit).Offset(offset).Order("name ASC").Find(&cities).Error
	if err != nil {
		return nil, 0, err
	}

	return cities, total, nil
}

func (r *GeolocationRepository) UpdateCity(city *domain.City) error {
	return r.db.Save(city).Error
}

func (r *GeolocationRepository) DeleteCity(id string) error {
	return r.db.Delete(&domain.City{}, "id = ?", id).Error
}

func (r *GeolocationRepository) SearchCities(query string, stateID string, limit, offset int) ([]*domain.City, int64, error) {
	var cities []*domain.City
	var total int64

	dbQuery := r.db.Model(&domain.City{}).Preload("State").Preload("State.Country")

	if stateID != "" {
		dbQuery = dbQuery.Where("state_id = ?", stateID)
	}

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where("LOWER(name) LIKE ? OR LOWER(zip_code) LIKE ?", 
			searchPattern, searchPattern)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get cities with pagination
	err := dbQuery.Limit(limit).Offset(offset).Order("name ASC").Find(&cities).Error
	if err != nil {
		return nil, 0, err
	}

	return cities, total, nil
}

// Utility operations
func (r *GeolocationRepository) GetCountriesWithStates(limit, offset int) ([]*domain.Country, int64, error) {
	var countries []*domain.Country
	var total int64

	// Get total count
	if err := r.db.Model(&domain.Country{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get countries with states
	err := r.db.Preload("States").Limit(limit).Offset(offset).Order("name ASC").Find(&countries).Error
	if err != nil {
		return nil, 0, err
	}

	return countries, total, nil
}

func (r *GeolocationRepository) GetStatesWithCities(countryID string, limit, offset int) ([]*domain.State, int64, error) {
	var states []*domain.State
	var total int64

	// Get total count
	if err := r.db.Model(&domain.State{}).Where("country_id = ?", countryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get states with cities
	err := r.db.Preload("Country").Preload("Cities").Where("country_id = ?", countryID).
		Limit(limit).Offset(offset).Order("name ASC").Find(&states).Error
	if err != nil {
		return nil, 0, err
	}

	return states, total, nil
}