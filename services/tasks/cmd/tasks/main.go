package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	taskshttp "github.com/student/tech-ip-sem2/services/tasks/internal/http"
	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
	"github.com/student/tech-ip-sem2/shared/logger"
)

func main() {
	log, err := logger.New("tasks", zap.InfoLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	tasksPort := os.Getenv("TASKS_PORT")
	if tasksPort == "" {
		tasksPort = "8082"
	}
	authGrpcAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGrpcAddr == "" {
		authGrpcAddr = "localhost:50051"
	}

	taskService := service.NewTaskService()
	authClient, err := grpcclient.NewAuthClient(authGrpcAddr)
	if err != nil {
		log.Fatal("failed to create auth gRPC client", zap.Error(err))
	}
	defer authClient.Close()

	router := taskshttp.NewRouter(taskService, authClient, log)

	server := &http.Server{
		Addr:         ":" + tasksPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Tasks service started", zap.String("port", tasksPort), zap.String("auth_grpc_addr", authGrpcAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", zap.Error(err))
		}
	}()

	<-done
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	}
	authClient.SetLogger(log)
}
