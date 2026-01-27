package http

import (
	"io"
	"net/http"
)

func (s *HTTPServer) test(w http.ResponseWriter, r *http.Request) {
	file, handler, _ := r.FormFile("file")
	s.resourceSrv.UploadResource(r.Context(), file, handler.Filename)

}

func v(r io.Reader) {

}
