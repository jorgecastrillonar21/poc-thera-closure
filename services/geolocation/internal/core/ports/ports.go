package ports

import (
	"theraclosure/geolocation-service/internal/core/domain"
)

// GeolocationRepository defines the interface for geolocation data access
type GeolocationRepository interface {
	// Country operations
	CreateCountry(country *domain.Country) error
	GetCountryByID(id string) (*domain.Country, error)
	GetCountryByCode(code string) (*domain.Country, error)
	GetAllCountries(limit, offset int) ([]*domain.Country, int64, error)
	UpdateCountry(country *domain.Country) error
	DeleteCountry(id string) error
	SearchCountries(query string, limit, offset int) ([]*domain.Country, int64, error)

	// State operations
	CreateState(state *domain.State) error
	GetStateByID(id string) (*domain.State, error)
	GetStatesByCountryID(countryID string, limit, offset int) ([]*domain.State, int64, error)
	UpdateState(state *domain.State) error
	DeleteState(id string) error
	SearchStates(query string, countryID string, limit, offset int) ([]*domain.State, int64, error)

	// City operations
	CreateCity(city *domain.City) error
	GetCityByID(id string) (*domain.City, error)
	GetCitiesByStateID(stateID string, limit, offset int) ([]*domain.City, int64, error)
	UpdateCity(city *domain.City) error
	DeleteCity(id string) error
	SearchCities(query string, stateID string, limit, offset int) ([]*domain.City, int64, error)

	// Utility operations
	GetCountriesWithStates(limit, offset int) ([]*domain.Country, int64, error)
	GetStatesWithCities(countryID string, limit, offset int) ([]*domain.State, int64, error)
}

// GeolocationService defines the business logic interface
type GeolocationService interface {
	// Country services
	CreateCountry(name, code, code2, region, currency string) (*domain.Country, error)
	GetCountry(id string) (*domain.Country, error)
	GetCountryByCode(code string) (*domain.Country, error)
	ListCountries(limit, offset int) ([]*domain.Country, int64, error)
	UpdateCountry(id, name, code, code2, region, currency string, active bool) (*domain.Country, error)
	DeleteCountry(id string) error
	SearchCountries(query string, limit, offset int) ([]*domain.Country, int64, error)

	// State services
	CreateState(countryID, name, code string) (*domain.State, error)
	GetState(id string) (*domain.State, error)
	ListStatesByCountry(countryID string, limit, offset int) ([]*domain.State, int64, error)
	UpdateState(id, name, code string, active bool) (*domain.State, error)
	DeleteState(id string) error
	SearchStates(query string, countryID string, limit, offset int) ([]*domain.State, int64, error)

	// City services
	CreateCity(stateID, name, zipCode string, latitude, longitude float64) (*domain.City, error)
	GetCity(id string) (*domain.City, error)
	ListCitiesByState(stateID string, limit, offset int) ([]*domain.City, int64, error)
	UpdateCity(id, name, zipCode string, latitude, longitude float64, active bool) (*domain.City, error)
	DeleteCity(id string) error
	SearchCities(query string, stateID string, limit, offset int) ([]*domain.City, int64, error)

	// Hierarchical services
	GetCompleteHierarchy(limit, offset int) ([]*domain.Country, int64, error)
	GetCountryHierarchy(countryID string) (*domain.Country, error)
	GetStateHierarchy(stateID string) (*domain.State, error)

	// Bulk operations
	BulkCreateCountries(countries []domain.Country) error
	BulkCreateStates(states []domain.State) error
	BulkCreateCities(cities []domain.City) error
}