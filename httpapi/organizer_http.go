package httpapi

import (
	"jazz-network-api/model"
	s "jazz-network-api/service"
	"net/http"
	"strings"
)

type OrganizerHandler struct {
	service *s.OrganizerService
}

func (h *OrganizerHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
	}

	organizer, err := h.service.Read(r.Context(), id)
	if err != nil {
		writeError(w, err)
	}

	writeJSON(w, http.StatusOK, prepareForJSON(organizer))
}

func (h *OrganizerHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	organizers, err := h.service.ReadAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	organizersJSON := make([]model.OrganizerJSON, len(organizers))
	for i, o := range organizers {
		organizersJSON[i] = prepareForJSON(o)
	}

	writeJSON(w, http.StatusOK, organizersJSON)
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
