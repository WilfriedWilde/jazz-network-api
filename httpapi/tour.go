package httpapi

import (
	"jazz-network-api/model"
	s "jazz-network-api/service"
	"net/http"
)

type TourHandler struct {
	Service *s.TourService
}

func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.TourInput](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	tour, err := h.Service.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tour)
}
