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

func (h *TourHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	tours, err := h.Service.ReadAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tours)
}

func (h *TourHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}

	tour, err := h.Service.Read(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tour)
}

func (h *TourHandler) Update(w http.ResponseWriter, r *http.Request) {
	tourID, err := parsePathID(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}

	update, err := decodeJSON[model.TourUpdate](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	tour, err := h.Service.Update(r.Context(), tourID, update)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tour)
}
