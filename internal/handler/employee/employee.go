package employee_handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimashiro/test_mediasoft/internal/middleware"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/dimashiro/test_mediasoft/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	employeeCreateURL = "/api/employees/create"
	employeeUpdateURL = "/api/employees/update"
	employeesURL      = "/api/employees"
	employeeDeleteURL = "/api/employees/delete"
)

type Handler struct {
	log   *zap.SugaredLogger
	uCase *usecase.Employee
}

func New(log *zap.SugaredLogger, uCase *usecase.Employee) Handler {
	return Handler{log: log, uCase: uCase}
}

func (h Handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, employeeCreateURL, middleware.Logging(h.log, h.Create))
	r.HandlerFunc(http.MethodGet, employeesURL, middleware.Logging(h.log, h.GetAll))
	r.HandlerFunc(http.MethodPut, employeeUpdateURL, middleware.Logging(h.log, h.Update))
	r.HandlerFunc(http.MethodDelete, employeeDeleteURL, middleware.Logging(h.log, h.Delete))
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// prepare dto to parse request
	ctx := r.Context()
	dto := &dto.CreateEmployee{}
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

	empl, err := h.uCase.CreateEmployee(ctx, dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't create employee: "+err.Error())
		http.Error(w, "can't create employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(empl)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't marshal employee: "+err.Error())
		http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't write json data: "+err.Error())
		http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	empls, err := h.uCase.GetAllEmployees(ctx)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't get employees: "+err.Error())
		http.Error(w, "can't get employees: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(empls)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't marshal employees: "+err.Error())
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

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dto := &dto.DeleteEmployee{}
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

	err = h.uCase.DeleteEmployee(ctx, dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't delete employee: "+err.Error())
		http.Error(w, "can't delete employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(`{"success": "ok"}`)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dto := &dto.UpdateEmployee{}
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

	err = h.uCase.UpdateEmployee(ctx, dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't update employee: "+err.Error())
		http.Error(w, "can't create employee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(`{"success": "ok"}`)
}

func (h Handler) validateReq(dto interface{}) error {
	//TODO add validation
	return nil
}
