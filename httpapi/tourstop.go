package httpapi

import (
	"jazz-network-api/model"
	"jazz-network-api/service"
	"net/http"
)

type TourStopHandler struct {
	Service *service.TourStopService
}

func (h *TourStopHandler) Update(w http.ResponseWriter, r *http.Request) {
	tourStopID, err := parsePathID(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}

	update, err := decodeJSON[model.TourStopUpdate](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	tourStop, err := h.Service.Update(r.Context(), tourStopID, update)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tourStop)
}
