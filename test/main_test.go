package test

import (
	ct "ClockTowerAPI/http"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := ct.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	assert.Equal(t, 200, w.Code)

	// for simplicity, let's unmarshal the body response to a map of string
	// in the future please use a struct for an application that has important things in it
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	msgRes, msgExists := response["message"]
	verRes, verExists := response["version"]

	//assert.Equal(t, `{"message":"ClockTower Backend","version":"0.0.1"}`, w.Body.JSON())
	assert.Nil(t, err, nil)
	assert.True(t, msgExists, true)
	assert.True(t, verExists, true)
	assert.Equal(t, msgRes, "ClockTowerAPI")
	assert.Equal(t, verRes, "0.0.1")
}
