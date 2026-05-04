package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type chatRequest struct {
	Message string `json:"message"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

func (srv *HTTPServer) ChatHandler(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseWithErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Message == "" {
		ResponseWithErr(w, http.StatusBadRequest, "message is required")
		return
	}

	reply, err := srv.geminiClient.Chat(r.Context(), req.Message)
	if err != nil {
		slog.Error("gemini chat error", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to get response from AI")
		return
	}

	ResponseWithJSON(w, http.StatusOK, chatResponse{Reply: reply})
}
