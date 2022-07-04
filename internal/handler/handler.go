package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dimashiro/test_mediasoft/config"
	department_handler "github.com/dimashiro/test_mediasoft/internal/handler/department"
	"github.com/dimashiro/test_mediasoft/internal/repository/department"
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
	repo := department.New(pool)
	departmentUCase := usecase.NewDepartment(log, repo)

	department_handler.New(log, departmentUCase).Register(router)
	return router, nil
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
