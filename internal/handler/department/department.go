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
	departmentCreateURL          = "/api/department/create"
	departmentUpdateURL          = "/api/department/update"
	departmentHierachyURL        = "/api/departments/hierarchy"
	departmentDeleteURL          = "/api/department/delete"
	emplInDepartmentURL          = "/api/department/:uuid/employees"
	emplInDepartmentHierarchyURL = "/api/department/:uuid/employees/all"
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
	r.HandlerFunc(http.MethodGet, departmentHierachyURL, h.Hierarchy)
	r.HandlerFunc(http.MethodDelete, departmentDeleteURL, h.Delete)
	r.HandlerFunc(http.MethodGet, emplInDepartmentURL, h.GetEmployees)
	r.HandlerFunc(http.MethodGet, emplInDepartmentHierarchyURL, h.GetEmployeesInHierarchy)
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

func (h Handler) Hierarchy(w http.ResponseWriter, r *http.Request) {
	dps, err := h.uCase.HierarchyDepartment()
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't get departments: "+err.Error())
		http.Error(w, "can't get departments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(dps)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't marshal departments: "+err.Error())
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
	// prepare dto to parse request
	dto := &dto.DeleteDepartment{}
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

	err = h.uCase.DeleteDepartment(dto)
	if err != nil {
		h.log.Errorw("ERROR", "ERROR", "can't delete department: "+err.Error())
		http.Error(w, "can't delete department: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(`{"success": "ok"}`)
}

func (h Handler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	emplUUID := params.ByName("uuid")
	if emplUUID == "" {
		h.log.Errorw("ERROR", "ERROR", "wrong uuid in req")
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	empls, err := h.uCase.GetEmployeesByDepartment(emplUUID)
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

func (h Handler) GetEmployeesInHierarchy(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	dpUUID := params.ByName("uuid")
	if dpUUID == "" {
		h.log.Errorw("ERROR", "ERROR", "wrong uuid in req")
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	empls, err := h.uCase.GetEmployeesInDepartmentHierarchy(dpUUID)
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

func (h Handler) validateReq(dto interface{}) error {
	//TODO add validation
	return nil
}
