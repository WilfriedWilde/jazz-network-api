package httpapi

import (
	"jazz-network-api/model"
	s "jazz-network-api/service"
	"net/http"
	"strings"
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

	writeJSON(w, http.StatusOK, prepareForJSON(organizer))
}

func (h *OrganizerHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	organizers, err := h.Service.ReadAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	organizersJSON := make([]model.OrganizerJSON, len(organizers))
	for i, o := range organizers {
		organizersJSON[i] = prepareForJSON(o)
	}

	writeJSON(w, http.StatusOK, organizersJSON)
}

func (h *OrganizerHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.OrganizerCreateInput](w, r)
	if err != nil {
		writeError(w, err)
	}

	organizer, err := h.Service.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, prepareForJSON(organizer))
}

func prepareForJSON(o model.Organizer) model.OrganizerJSON {
	return model.OrganizerJSON{
		ID:             o.ID,
		Type:           o.Type,
		Name:           o.Name,
		Country:        o.Country,
		City:           o.City,
		Address:        o.Address,
		Websites:       strings.Split(o.Websites, ";"),
		Emails:         strings.Split(o.Emails, ";"),
		PhoneNumbers:   strings.Split(o.PhoneNumbers, ";"),
		ContactPersons: strings.Split(o.ContactPersons, ";"),
		Schedule:       o.Schedule,
		Notes:          o.Notes,
		Latitude:       o.Latitude,
		Longitude:      o.Longitude,
	}
}
