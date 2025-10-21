package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"theraclosure/geolocation-service/internal/core/domain"
)

// Country handlers
func (s *Server) createCountry(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Code2    string `json:"code2" binding:"required"`
		Region   string `json:"region"`
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	country, err := s.service.CreateCountry(req.Name, req.Code, req.Code2, req.Region, req.Currency)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Country created successfully",
		"country": country,
	})
}

func (s *Server) getCountry(c *gin.Context) {
	id := c.Param("id")

	country, err := s.service.GetCountry(id)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	if country == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, country)
}

func (s *Server) listCountries(c *gin.Context) {
	limit, offset := s.getPaginationParams(c)

	countries, total, err := s.service.ListCountries(limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (s *Server) updateCountry(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name     string `json:"name"`
		Code     string `json:"code"`
		Code2    string `json:"code2"`
		Region   string `json:"region"`
		Currency string `json:"currency"`
		Active   *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	country, err := s.service.UpdateCountry(id, req.Name, req.Code, req.Code2, req.Region, req.Currency, active)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Country updated successfully",
		"country": country,
	})
}

func (s *Server) deleteCountry(c *gin.Context) {
	id := c.Param("id")

	err := s.service.DeleteCountry(id)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Country deleted successfully",
	})
}

func (s *Server) searchCountries(c *gin.Context) {
	query := c.Query("q")
	limit, offset := s.getPaginationParams(c)

	countries, total, err := s.service.SearchCountries(query, limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
		"query":     query,
	})
}

// State handlers
func (s *Server) createState(c *gin.Context) {
	var req struct {
		CountryID string `json:"country_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	state, err := s.service.CreateState(req.CountryID, req.Name, req.Code)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "State created successfully",
		"state":   state,
	})
}

func (s *Server) getState(c *gin.Context) {
	id := c.Param("id")

	state, err := s.service.GetState(id)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	if state == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "State not found"})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (s *Server) listStatesByCountry(c *gin.Context) {
	countryID := c.Param("id")
	limit, offset := s.getPaginationParams(c)

	states, total, err := s.service.ListStatesByCountry(countryID, limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"states":     states,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"country_id": countryID,
	})
}

func (s *Server) updateState(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Active *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	state, err := s.service.UpdateState(id, req.Name, req.Code, active)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "State updated successfully",
		"state":   state,
	})
}

func (s *Server) deleteState(c *gin.Context) {
	id := c.Param("id")

	err := s.service.DeleteState(id)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "State deleted successfully",
	})
}

func (s *Server) searchStates(c *gin.Context) {
	query := c.Query("q")
	countryID := c.Query("country_id")
	limit, offset := s.getPaginationParams(c)

	states, total, err := s.service.SearchStates(query, countryID, limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"states":     states,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"query":      query,
		"country_id": countryID,
	})
}

// City handlers
func (s *Server) createCity(c *gin.Context) {
	var req struct {
		StateID   string  `json:"state_id" binding:"required"`
		Name      string  `json:"name" binding:"required"`
		ZipCode   string  `json:"zip_code"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	city, err := s.service.CreateCity(req.StateID, req.Name, req.ZipCode, req.Latitude, req.Longitude)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "City created successfully",
		"city":    city,
	})
}

func (s *Server) getCity(c *gin.Context) {
	id := c.Param("id")

	city, err := s.service.GetCity(id)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	if city == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		return
	}

	c.JSON(http.StatusOK, city)
}

func (s *Server) listCitiesByState(c *gin.Context) {
	stateID := c.Param("id")
	limit, offset := s.getPaginationParams(c)

	cities, total, err := s.service.ListCitiesByState(stateID, limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cities":   cities,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"state_id": stateID,
	})
}

func (s *Server) updateCity(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string   `json:"name"`
		ZipCode   string   `json:"zip_code"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Active    *bool    `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	latitude := 0.0
	if req.Latitude != nil {
		latitude = *req.Latitude
	}

	longitude := 0.0
	if req.Longitude != nil {
		longitude = *req.Longitude
	}

	city, err := s.service.UpdateCity(id, req.Name, req.ZipCode, latitude, longitude, active)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "City updated successfully",
		"city":    city,
	})
}

func (s *Server) deleteCity(c *gin.Context) {
	id := c.Param("id")

	err := s.service.DeleteCity(id)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "City deleted successfully",
	})
}

func (s *Server) searchCities(c *gin.Context) {
	query := c.Query("q")
	stateID := c.Query("state_id")
	limit, offset := s.getPaginationParams(c)

	cities, total, err := s.service.SearchCities(query, stateID, limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cities":   cities,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"query":    query,
		"state_id": stateID,
	})
}

// Hierarchy handlers
func (s *Server) getCompleteHierarchy(c *gin.Context) {
	limit, offset := s.getPaginationParams(c)

	countries, total, err := s.service.GetCompleteHierarchy(limit, offset)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hierarchy": countries,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (s *Server) getCountryHierarchy(c *gin.Context) {
	countryID := c.Param("id")

	country, err := s.service.GetCountryHierarchy(countryID)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	if country == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, country)
}

func (s *Server) getStateHierarchy(c *gin.Context) {
	stateID := c.Param("id")

	state, err := s.service.GetStateHierarchy(stateID)
	if err != nil {
		s.handleError(c, err, http.StatusInternalServerError)
		return
	}

	if state == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "State not found"})
		return
	}

	c.JSON(http.StatusOK, state)
}

// Bulk operation handlers
func (s *Server) bulkCreateCountries(c *gin.Context) {
	var req struct {
		Countries []domain.Country `json:"countries" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := s.service.BulkCreateCountries(req.Countries)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Countries created successfully",
		"count":   len(req.Countries),
	})
}

func (s *Server) bulkCreateStates(c *gin.Context) {
	var req struct {
		States []domain.State `json:"states" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := s.service.BulkCreateStates(req.States)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "States created successfully",
		"count":   len(req.States),
	})
}

func (s *Server) bulkCreateCities(c *gin.Context) {
	var req struct {
		Cities []domain.City `json:"cities" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := s.service.BulkCreateCities(req.Cities)
	if err != nil {
		s.handleError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cities created successfully",
		"count":   len(req.Cities),
	})
}