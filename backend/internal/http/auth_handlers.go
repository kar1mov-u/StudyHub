package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *HTTPServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, err := s.authSrv.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to create the JWT")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
