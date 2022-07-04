package department_handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/dimashiro/test_mediasoft/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	departmentCreateURL = "/api/departments/create"
	departmentUpdateURL = "/api/departments/update"
)

type Handler struct {
	log   *zap.SugaredLogger
	uCase *usecase.Department
}

func New(log *zap.SugaredLogger, uCase *usecase.Department) Handler {
	return Handler{log: log, uCase: uCase}
}

func (h Handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, departmentCreateURL, h.Create)
	r.HandlerFunc(http.MethodPost, departmentUpdateURL, h.Update)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// prepare dto to parse request
	dto := &dto.CreateDepartment{}
	// parse req body to dto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't parse req: "+err.Error())
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check that request valid
	err = h.validateReq(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't parse req: "+err.Error())
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	dp, err := h.uCase.CreateDepartment(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't create department: "+err.Error())
		http.Error(w, "can't create department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(dp)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't marshal department: "+err.Error())
		http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't write json data: "+err.Error())
		http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	// prepare dto to parse request
	dto := &dto.UpdateDepartment{}
	// parse req body to dto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't parse req: "+err.Error())
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check that request valid
	err = h.validateReq(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't parse req: "+err.Error())
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.uCase.UpdateDepartment(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't update department: "+err.Error())
		http.Error(w, "can't create department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(`{"success": "ok"}`)
}

func (h Handler) validateReq(dto interface{}) error {
	//TODO add validation
	return nil
}
