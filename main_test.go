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

	// Здесь только тестируемые значения count
	// Ожидаемые значения вычисляются автоматически в процессе проверки
	requests := []int{0, 1, 2, 100}

	for _, v := range requests {
		c :="moscow"

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", v, c), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}

		// Ожидаемые значения вычисляются автоматически
		assert.Len(t, cafes, min(v, len(cafeList[c])))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	// Здесь только тестируемые значения search
	// Ожидаемые значения вычисляются автоматически и 
	// подставляются сразу после объявления структуры с помощью цикла for
	requests := []struct{
		search string
		wantCount int
	}{
		{search: "фасоль"},
		{search: "кофе"},
		{search: "вилка"},
	}

	// Здесь вычисляются и подставляются ожидаемые значения
	for i, s := range requests {
		for _, c := range cafeList["moscow"] {
			if strings.Contains(strings.ToLower(c), strings.ToLower(s.search)) {
				requests[i].wantCount++
			}
		}
	}

	// Здесь уже запускаются проверки
	for _, s := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", s.search), nil)

		handler.ServeHTTP(response, req)
		
		require.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if cafes[0] == "" {
			cafes = []string{}
		}
		assert.Len(t, cafes, s.wantCount)
	}
}