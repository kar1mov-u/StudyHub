package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type chatRequest struct {
	Message string `json:"message"`
}

type chatSource struct {
	Source string `json:"source"`
	Page   *int   `json:"page,omitempty"`
}

type chatResponse struct {
	Reply   string       `json:"reply"`
	Sources []chatSource `json:"sources,omitempty"`
}

type ragChatResponse struct {
	Reply   string       `json:"reply"`
	Sources []chatSource `json:"sources"`
}

var ragHTTPClient = &http.Client{Timeout: 30 * time.Second}

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

	if srv.ragServiceURL != "" {
		if reply, sources, err := callRAGService(srv.ragServiceURL, req.Message); err == nil {
			ResponseWithJSON(w, http.StatusOK, chatResponse{Reply: reply, Sources: sources})
			return
		} else {
			slog.Warn("rag service unavailable, falling back to gemini", "err", err)
		}
	}

	reply, err := srv.geminiClient.Chat(r.Context(), req.Message)
	if err != nil {
		slog.Error("gemini chat error", "err", err)
		ResponseWithErr(w, http.StatusInternalServerError, "failed to get response from AI")
		return
	}
	ResponseWithJSON(w, http.StatusOK, chatResponse{Reply: reply})
}

func callRAGService(baseURL, message string) (string, []chatSource, error) {
	body, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		return "", nil, err
	}

	resp, err := ragHTTPClient.Post(baseURL+"/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", nil, fmt.Errorf("rag service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("rag service returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var ragResp ragChatResponse
	if err := json.Unmarshal(data, &ragResp); err != nil {
		return "", nil, err
	}
	if ragResp.Reply == "" {
		return "", nil, fmt.Errorf("empty reply from rag service")
	}
	return ragResp.Reply, ragResp.Sources, nil
}
