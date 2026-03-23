package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	"github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers"
	authMiddleware "github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers/middleware"
	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
	"github.com/student/tech-ip-sem2/shared/metrics"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func NewRouter(taskService *service.TaskService, authClient *grpcclient.AuthClient, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Метрики endpoint (без авторизации)
	mux.Handle("GET /metrics", promhttp.Handler())

	protected := authMiddleware.AuthMiddleware(authClient)

	mux.Handle("POST /v1/tasks", protected(http.HandlerFunc(handlers.CreateTaskHandler(taskService))))
	mux.Handle("GET /v1/tasks", protected(http.HandlerFunc(handlers.GetTasksHandler(taskService))))
	mux.Handle("GET /v1/tasks/{id}", protected(http.HandlerFunc(handlers.GetTaskHandler(taskService))))
	mux.Handle("PATCH /v1/tasks/{id}", protected(http.HandlerFunc(handlers.UpdateTaskHandler(taskService))))
	mux.Handle("DELETE /v1/tasks/{id}", protected(http.HandlerFunc(handlers.DeleteTaskHandler(taskService))))

	// Применяем middleware: request-id -> access log -> metrics
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.HTTPAccessLogMiddleware(logger)(handler)
	handler = metrics.HTTPMetricsMiddleware(handler) // <-- добавляем сбор метрик

	return handler
}
