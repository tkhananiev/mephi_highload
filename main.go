// @title Highload User Service
// @version 1.0
// @description REST API for user management
// @host localhost:8080
// @BasePath /

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "go-microservice/docs"
	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/services"
	"go-microservice/utils"
)

func main() {
	port := getenv("PORT", "8080")

	userSvc := services.NewUserService()

	audit := handlers.NewAudit(10_000)
	notify := handlers.NewNotifier(10_000)

	integSvc, integErr := services.NewIntegrationService(
		getenv("MINIO_ENDPOINT", "minio:9000"),
		getenv("MINIO_ACCESS_KEY", "minioadmin"),
		getenv("MINIO_SECRET_KEY", "minioadmin"),
		getenv("MINIO_BUCKET", "audit-bucket"),
		false,
	)
	if integErr != nil {
		utils.Logger.Printf("minio init error: %v (integration endpoints disabled)", integErr)
	}

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//rl := utils.NewRateLimiter(rate.Limit(1000), 5000)
	//r.Use(rl.Middleware)
	r.Use(metrics.MetricsMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	uh := handlers.NewUserHandler(userSvc, audit, notify)
	uh.Register(r)

	if integErr == nil {
		ih := handlers.NewIntegrationHandler(integSvc)
		r.HandleFunc("/api/integrations/audit/upload", ih.UploadAudit).Methods(http.MethodPost)
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	utils.Logger.Printf("listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Logger.Fatalf("server error: %v", err)
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
