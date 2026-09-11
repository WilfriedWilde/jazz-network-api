package httpapi

import (
	"jazz-network-api/model"
	s "jazz-network-api/service"
	"net/http"
)

type OrganizerHandler struct {
	Service *s.OrganizerService
}

func (h *OrganizerHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}

	organizer, err := h.Service.Read(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, organizer.ToOrganizerJSON())
}

func (h *OrganizerHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	organizers, err := h.Service.ReadAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	organizersJSON := make([]model.OrganizerJSON, len(organizers))
	for i, o := range organizers {
		organizersJSON[i] = o.ToOrganizerJSON()
	}

	writeJSON(w, http.StatusOK, organizersJSON)
}

func (h *OrganizerHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.OrganizerInput](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	organizer, err := h.Service.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, organizer.ToOrganizerJSON())
}
