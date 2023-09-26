package main

import (
	"Mahathirrr/belajar-golang-restful-api/config"
	"Mahathirrr/belajar-golang-restful-api/controller"
	"Mahathirrr/belajar-golang-restful-api/helper"
	"Mahathirrr/belajar-golang-restful-api/repository"
	"Mahathirrr/belajar-golang-restful-api/service"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {

	db := config.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := config.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
