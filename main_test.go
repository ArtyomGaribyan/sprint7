package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string {
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct{
		count, want int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 100},
	}
	c :="moscow"

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", v.count, c), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}

		assert.Len(t, cafes, min(v.want, len(cafeList[c])))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct{
		search string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)

		handler.ServeHTTP(response, req)
		
		require.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}
		assert.Len(t, cafes, v.wantCount)
		
		for _, cafe := range cafes {
			assert.Contains(t, strings.ToLower(cafe), v.search)
		}
	}
}