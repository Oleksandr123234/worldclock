package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testZones() []zoneInfo {
	return []zoneInfo{
		{Name: "UTC", OffsetMins: 0},
		{Name: "Europe/Kyiv", OffsetMins: 120},
	}
}

func TestHealthzEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(testZones())

	req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}

func TestIndexEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(testZones())

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "UTC")
	assert.Contains(t, rr.Body.String(), "Europe/Kyiv")
}

func TestLoadZones(t *testing.T) {
	t.Setenv("TIMEZONES", "UTC,Europe/Kyiv")

	zones := loadZones()

	assert.Len(t, zones, 2)
	assert.Equal(t, "UTC", zones[0].Name)
	assert.Equal(t, "Europe/Kyiv", zones[1].Name)
}
