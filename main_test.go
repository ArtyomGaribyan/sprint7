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

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
		require.Equal(t, http.StatusBadRequest, response.Code)
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

		assert.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []int{0, 1, 2, 100}

	for _, v := range requests {
		c :="moscow"

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", v, c), nil)

		handler.ServeHTTP(response, req)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}
		if v <= len(cafeList[c]) {
			assert.Equal(t, v, len(cafes))
		} else {
			assert.Equal(t, len(cafeList[c]), len(cafes))
		}
		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct{
		search string
		wantCount int
	}{
		{search: "фасоль"},
		{search: "кофе"},
		{search: "вилка"},
	}

	for i, s := range requests {
		for _, c := range cafeList["moscow"] {
			if strings.Contains(strings.ToLower(c), strings.ToLower(s.search)) {
				requests[i].wantCount++
			}
		}
	}

	fmt.Println(requests)

	for _, s := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", s.search), nil)

		handler.ServeHTTP(response, req)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}
		assert.Equal(t, s.wantCount, len(cafes))
		require.Equal(t, http.StatusOK, response.Code)
	}
}