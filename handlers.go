package main

import (
	"encoding/json"
	"jazz-network-api/repository"
	"log"
	"net/http"
)

func (s *Server) HandleSelectAllOrganizers(w http.ResponseWriter, r *http.Request) {
	organizers, err := repository.SelectAllOrganizers(r.Context(), s.db)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(organizers); err != nil {
		log.Printf("error ecnoding organizers: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
