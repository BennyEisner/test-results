package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/BennyEisner/test-results/internal/domain"
)

type BuildHandler struct {
	svc domain.BuildService
}

func NewBuildHandler(svc domain.BuildService) *BuildHandler {
	return &BuildHandler{svc: svc}
}

func (h *BuildHandler) GetBuildByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	ctx := r.Context()
	build, err := h.svc.GetBuildByID(ctx, id)
	if err != nil {
		switch err {
		case domain.ErrBuildNotFound:
			respondWithError(w, http.StatusNotFound, "build not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, build)
}

func (h *BuildHandler) GetBuildsByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	ctx := r.Context()
	builds, err := h.svc.GetBuildsByProjectID(ctx, projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get builds")
		return
	}
	respondWithJSON(w, http.StatusOK, builds)
}

func (h *BuildHandler) GetBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("suiteId")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test suite ID")
		return
	}
	ctx := r.Context()
	builds, err := h.svc.GetBuildsByTestSuiteID(ctx, suiteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get builds")
		return
	}
	respondWithJSON(w, http.StatusOK, builds)
}

func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var req domain.Build
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.CreatedAt = time.Now()
	ctx := r.Context()
	build, err := h.svc.CreateBuild(ctx, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create build")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, build)
}

func (h *BuildHandler) UpdateBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	var req domain.Build
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	build, err := h.svc.UpdateBuild(ctx, id, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update build")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, build)
}

func (h *BuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	ctx := r.Context()
	if err := h.svc.DeleteBuild(ctx, id); err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete build")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "build deleted successfully"})
}
