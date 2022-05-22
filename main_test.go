package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GET_File(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/file/app/1.bin?filterByName=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test binary file")
}

func Test_GET_Directory(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/file/app?orderBy=lastModified&orderByDirection=Descending&filterByName=g", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "{\"files\":[\".git\",\"main_test.go\",\"main.go\",\"go.mod\",\"go.sum\",\"截圖 2022-05-19 上午3.52.50.png\",\"照片.jpg\",\"截圖 2022-05-11 下午12.03.09.png\"],\"isDirectory\":true}")
}
