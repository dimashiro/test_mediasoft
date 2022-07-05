package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dimashiro/test_mediasoft/config"
	department_handler "github.com/dimashiro/test_mediasoft/internal/handler/department"
	employee_handler "github.com/dimashiro/test_mediasoft/internal/handler/employee"
	"github.com/dimashiro/test_mediasoft/internal/repository/department"
	"github.com/dimashiro/test_mediasoft/internal/repository/employee"
	"github.com/dimashiro/test_mediasoft/internal/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func NewRouter(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config) (*httprouter.Router, error) {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/heartbeat", Heartbeat)

	pool, err := pgxpool.Connect(ctx, cfg.GetDBConnString())
	if err != nil {
		return nil, fmt.Errorf("can't create pg pool: %s", err.Error())
	}
	rDptm := department.New(pool)
	departmentUCase := usecase.NewDepartment(log, rDptm)
	rEmpl := employee.New(pool)
	employeeUCase := usecase.NewEmployee(log, rEmpl, rDptm)

	employee_handler.New(log, employeeUCase).Register(router)
	department_handler.New(log, departmentUCase).Register(router)
	return router, nil
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
