package services

import (
	"errors"
	"strings"

	"theraclosure/geolocation-service/internal/core/domain"
	"theraclosure/geolocation-service/internal/core/ports"
)

type GeolocationService struct {
	repo ports.GeolocationRepository
}

func NewGeolocationService(repo ports.GeolocationRepository) ports.GeolocationService {
	return &GeolocationService{
		repo: repo,
	}
}

// Country Services
func (s *GeolocationService) CreateCountry(name, code, code2, region, currency string) (*domain.Country, error) {
	// Validate input
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("country name is required")
	}
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("country code is required")
	}
	if strings.TrimSpace(code2) == "" {
		return nil, errors.New("country code2 is required")
	}

	// Check if country with same code already exists
	existing, _ := s.repo.GetCountryByCode(code)
	if existing != nil {
		return nil, errors.New("country with this code already exists")
	}

	country := &domain.Country{
		Name:     strings.TrimSpace(name),
		Code:     strings.ToUpper(strings.TrimSpace(code)),
		Code2:    strings.ToUpper(strings.TrimSpace(code2)),
		Region:   strings.TrimSpace(region),
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
		Active:   true,
	}

	if !country.IsValid() {
		return nil, errors.New("invalid country data")
	}

	err := s.repo.CreateCountry(country)
	if err != nil {
		return nil, err
	}

	return country, nil
}

func (s *GeolocationService) GetCountry(id string) (*domain.Country, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("country ID is required")
	}

	return s.repo.GetCountryByID(id)
}

func (s *GeolocationService) GetCountryByCode(code string) (*domain.Country, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("country code is required")
	}

	return s.repo.GetCountryByCode(strings.ToUpper(code))
}

func (s *GeolocationService) ListCountries(limit, offset int) ([]*domain.Country, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetAllCountries(limit, offset)
}

func (s *GeolocationService) UpdateCountry(id, name, code, code2, region, currency string, active bool) (*domain.Country, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("country ID is required")
	}

	country, err := s.repo.GetCountryByID(id)
	if err != nil {
		return nil, err
	}

	if country == nil {
		return nil, errors.New("country not found")
	}

	// Update fields if provided
	if strings.TrimSpace(name) != "" {
		country.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(code) != "" {
		country.Code = strings.ToUpper(strings.TrimSpace(code))
	}
	if strings.TrimSpace(code2) != "" {
		country.Code2 = strings.ToUpper(strings.TrimSpace(code2))
	}
	if strings.TrimSpace(region) != "" {
		country.Region = strings.TrimSpace(region)
	}
	if strings.TrimSpace(currency) != "" {
		country.Currency = strings.ToUpper(strings.TrimSpace(currency))
	}
	country.Active = active

	if !country.IsValid() {
		return nil, errors.New("invalid country data")
	}

	err = s.repo.UpdateCountry(country)
	if err != nil {
		return nil, err
	}

	return country, nil
}

func (s *GeolocationService) DeleteCountry(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("country ID is required")
	}

	country, err := s.repo.GetCountryByID(id)
	if err != nil {
		return err
	}

	if country == nil {
		return errors.New("country not found")
	}

	return s.repo.DeleteCountry(id)
}

func (s *GeolocationService) SearchCountries(query string, limit, offset int) ([]*domain.Country, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.SearchCountries(strings.TrimSpace(query), limit, offset)
}

// State Services
func (s *GeolocationService) CreateState(countryID, name, code string) (*domain.State, error) {
	if strings.TrimSpace(countryID) == "" {
		return nil, errors.New("country ID is required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("state name is required")
	}

	// Verify country exists
	country, err := s.repo.GetCountryByID(countryID)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, errors.New("country not found")
	}

	state := &domain.State{
		CountryID: countryID,
		Name:      strings.TrimSpace(name),
		Code:      strings.ToUpper(strings.TrimSpace(code)),
		Active:    true,
	}

	if !state.IsValid() {
		return nil, errors.New("invalid state data")
	}

	err = s.repo.CreateState(state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *GeolocationService) GetState(id string) (*domain.State, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("state ID is required")
	}

	return s.repo.GetStateByID(id)
}

func (s *GeolocationService) ListStatesByCountry(countryID string, limit, offset int) ([]*domain.State, int64, error) {
	if strings.TrimSpace(countryID) == "" {
		return nil, 0, errors.New("country ID is required")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetStatesByCountryID(countryID, limit, offset)
}

func (s *GeolocationService) UpdateState(id, name, code string, active bool) (*domain.State, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("state ID is required")
	}

	state, err := s.repo.GetStateByID(id)
	if err != nil {
		return nil, err
	}

	if state == nil {
		return nil, errors.New("state not found")
	}

	// Update fields if provided
	if strings.TrimSpace(name) != "" {
		state.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(code) != "" {
		state.Code = strings.ToUpper(strings.TrimSpace(code))
	}
	state.Active = active

	if !state.IsValid() {
		return nil, errors.New("invalid state data")
	}

	err = s.repo.UpdateState(state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *GeolocationService) DeleteState(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("state ID is required")
	}

	state, err := s.repo.GetStateByID(id)
	if err != nil {
		return err
	}

	if state == nil {
		return errors.New("state not found")
	}

	return s.repo.DeleteState(id)
}

func (s *GeolocationService) SearchStates(query string, countryID string, limit, offset int) ([]*domain.State, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.SearchStates(strings.TrimSpace(query), countryID, limit, offset)
}

// City Services
func (s *GeolocationService) CreateCity(stateID, name, zipCode string, latitude, longitude float64) (*domain.City, error) {
	if strings.TrimSpace(stateID) == "" {
		return nil, errors.New("state ID is required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("city name is required")
	}

	// Verify state exists
	state, err := s.repo.GetStateByID(stateID)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errors.New("state not found")
	}

	city := &domain.City{
		StateID:   stateID,
		Name:      strings.TrimSpace(name),
		ZipCode:   strings.TrimSpace(zipCode),
		Latitude:  latitude,
		Longitude: longitude,
		Active:    true,
	}

	if !city.IsValid() {
		return nil, errors.New("invalid city data")
	}

	err = s.repo.CreateCity(city)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func (s *GeolocationService) GetCity(id string) (*domain.City, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("city ID is required")
	}

	return s.repo.GetCityByID(id)
}

func (s *GeolocationService) ListCitiesByState(stateID string, limit, offset int) ([]*domain.City, int64, error) {
	if strings.TrimSpace(stateID) == "" {
		return nil, 0, errors.New("state ID is required")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetCitiesByStateID(stateID, limit, offset)
}

func (s *GeolocationService) UpdateCity(id, name, zipCode string, latitude, longitude float64, active bool) (*domain.City, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("city ID is required")
	}

	city, err := s.repo.GetCityByID(id)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, errors.New("city not found")
	}

	// Update fields if provided
	if strings.TrimSpace(name) != "" {
		city.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(zipCode) != "" {
		city.ZipCode = strings.TrimSpace(zipCode)
	}
	if latitude != 0 {
		city.Latitude = latitude
	}
	if longitude != 0 {
		city.Longitude = longitude
	}
	city.Active = active

	if !city.IsValid() {
		return nil, errors.New("invalid city data")
	}

	err = s.repo.UpdateCity(city)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func (s *GeolocationService) DeleteCity(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("city ID is required")
	}

	city, err := s.repo.GetCityByID(id)
	if err != nil {
		return err
	}

	if city == nil {
		return errors.New("city not found")
	}

	return s.repo.DeleteCity(id)
}

func (s *GeolocationService) SearchCities(query string, stateID string, limit, offset int) ([]*domain.City, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.SearchCities(strings.TrimSpace(query), stateID, limit, offset)
}

// Hierarchical Services
func (s *GeolocationService) GetCompleteHierarchy(limit, offset int) ([]*domain.Country, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetCountriesWithStates(limit, offset)
}

func (s *GeolocationService) GetCountryHierarchy(countryID string) (*domain.Country, error) {
	if strings.TrimSpace(countryID) == "" {
		return nil, errors.New("country ID is required")
	}

	return s.repo.GetCountryByID(countryID)
}

func (s *GeolocationService) GetStateHierarchy(stateID string) (*domain.State, error) {
	if strings.TrimSpace(stateID) == "" {
		return nil, errors.New("state ID is required")
	}

	return s.repo.GetStateByID(stateID)
}

// Bulk Operations
func (s *GeolocationService) BulkCreateCountries(countries []domain.Country) error {
	if len(countries) == 0 {
		return errors.New("no countries provided")
	}

	for i := range countries {
		countries[i].Active = true
		if !countries[i].IsValid() {
			return errors.New("invalid country data at index " + string(rune(i)))
		}
	}

	for i := range countries {
		err := s.repo.CreateCountry(&countries[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GeolocationService) BulkCreateStates(states []domain.State) error {
	if len(states) == 0 {
		return errors.New("no states provided")
	}

	for i := range states {
		states[i].Active = true
		if !states[i].IsValid() {
			return errors.New("invalid state data at index " + string(rune(i)))
		}
	}

	for i := range states {
		err := s.repo.CreateState(&states[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GeolocationService) BulkCreateCities(cities []domain.City) error {
	if len(cities) == 0 {
		return errors.New("no cities provided")
	}

	for i := range cities {
		cities[i].Active = true
		if !cities[i].IsValid() {
			return errors.New("invalid city data at index " + string(rune(i)))
		}
	}

	for i := range cities {
		err := s.repo.CreateCity(&cities[i])
		if err != nil {
			return err
		}
	}

	return nil
}