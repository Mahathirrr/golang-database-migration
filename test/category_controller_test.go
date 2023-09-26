package test

import (
	"Mahathirrr/belajar-golang-restful-api/config"
	"Mahathirrr/belajar-golang-restful-api/controller"
	"Mahathirrr/belajar-golang-restful-api/helper"
	"Mahathirrr/belajar-golang-restful-api/middleware"
	"Mahathirrr/belajar-golang-restful-api/model/domain"
	"Mahathirrr/belajar-golang-restful-api/repository"
	"Mahathirrr/belajar-golang-restful-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
	- Integration test --> hit ke web APInya langsung (controller)
*/

// Ada baiknya kita pisah testing db dengan main database dari project kita, karena tidak ingin menggangu data yang ada di main database
// karena nanti kita testnya ada masukin data
func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/belajar_golang_restful_api")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := config.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// akan membersihkan data table terlebih dahulu, sebelum melakukan testing yang akan memasukkan data ke table database
	// supaya tidak pusing lagi testing dan lebih clear,jelas sesuai dengan data per testingnya saja yang ditampilkan
	truncateCategory(db)

	requestBody := strings.NewReader(`{"name" : "Gadget"}"`)
	request := httptest.NewRequest(http.MethodPost, "https://localhost:3000/api/categories", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Result().StatusCode)

	body, _ := io.ReadAll(recorder.Result().Body)
	// simpan json di golang dalam bentuk map
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// responseBody["code"].(float64) mengonversi nilai tersebut menjadi float64 karena Go akan menganggap angka dalam JSON sebagai float64 secara default.

	// expected --> int(200), maka untuk Actual kita perlu melalukan konversi dan parsing karena valuenya berupa interface{}
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	// --> Bad Request
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodPut, "https://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// response.status --> 200 OK (gabungan antara status dan statusCode)
	// berbeda dengan apabila mengambil statusnya pada responseBody
	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPut, "https://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "https://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// response.status --> 200 OK ( gabungan antara status dan statusCode)
	// berbeda dengan apabila mengambil statusnya pada responseBody
	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "https://localhost:3000/api/categories/404", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "https://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "https://localhost:3000/api/categories/404", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestListCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "https://localhost:3000/api/categories", nil)
	request.Header.Set("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	fmt.Println(responseBody)

	// responbody yang berupa map[string]interface{} --> valuenya kan interface{} --> maka interface{} tersebut tidak bisa dikonversi ke map
	// karena yang bisa dikonversi ke map adalah data berupa slice, oleh karena itu perlu di konversi terlebih dahulu menjadi []interface{}
	/*
		var categories = responseBody["data"].([]map[string]interface{})

		assert.Equal(t, category1.Id, categories[0]["id"])
		assert.Equal(t, category1.Name, categories[0]["name"])
		assert.Equal(t, category2.Name, categories[1]["name"])
		assert.Equal(t, category2.Name, categories[1]["name"])
	*/

	var categories = responseBody["data"].([]interface{})
	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))

	assert.Equal(t, category1.Name, categoryResponse1["name"])
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "https://localhost:3000/api/categories", nil)
	request.Header.Set("X-API-KEY", "SALAH")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	fmt.Println(response.Status)
	fmt.Println(responseBody["status"])

	assert.Equal(t, 401, response.StatusCode)
	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}
