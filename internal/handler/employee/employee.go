package employee_handler

import (
	"encoding/json"
	"net/http"

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
	r.HandlerFunc(http.MethodPost, employeeCreateURL, h.Create)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	// prepare dto to parse request
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

	empl, err := h.uCase.CreateEmployee(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't create department: "+err.Error())
		http.Error(w, "can't create department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(empl)
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

func (h Handler) validateReq(dto interface{}) error {
	//TODO add validation
	return nil
}
