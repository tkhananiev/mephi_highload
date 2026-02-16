package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-microservice/services"
)

type IntegrationHandler struct {
	svc *services.IntegrationService
}

func NewIntegrationHandler(svc *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{svc: svc}
}

type uploadAuditResponse struct {
	Object   string `json:"object"`
	Location string `json:"location"`
}

// UploadAudit godoc
// @Summary Upload audit snapshot to MinIO
// @Description Creates bucket if missing and uploads a small audit snapshot object to MinIO
// @Tags integrations
// @Produce json
// @Success 200 {object} uploadAuditResponse
// @Failure 502 {string} string
// @Router /api/integrations/audit/upload [post]
func (h *IntegrationHandler) UploadAudit(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.svc.EnsureBucket(ctx); err != nil {
		http.Error(w, "minio bucket error", http.StatusBadGateway)
		return
	}

	obj := services.DefaultAuditObjectName()
	content := []byte("audit snapshot: service is running")

	loc, err := h.svc.UploadText(ctx, obj, content)
	if err != nil {
		http.Error(w, "minio upload error", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(uploadAuditResponse{
		Object:   obj,
		Location: loc,
	})
}
