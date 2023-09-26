package exception

import (
	"Mahathirrr/belajar-golang-restful-api/helper"
	"Mahathirrr/belajar-golang-restful-api/model/web"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {

	// router akan menjalankan function tersebut, ditambahkan kondisi if & return supaya, hanya 1 error saja yang menjadi responnya,
	// ketika return, program langsung berhenti, tidak melanjutkan kebawah
	if notFoundError(w, r, err) {
		return
	}
	if validationErrors(w, r, err) {
		return
	}
	internalServerError(w, r, err)
}

func validationErrors(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}
		helper.WriteToResponseBody(w, webResponse)
		return true
	} else {
		return false
	}
}

func notFoundError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}
		helper.WriteToResponseBody(w, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}
	helper.WriteToResponseBody(w, webResponse)
}
