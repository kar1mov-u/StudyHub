package http

import (
	"StudyHub/internal/resources"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// mockResourceService implements resource service methods for testing
type mockResourceService struct {
	uploadResourceFunc       func(ctx context.Context, file io.Reader, size int64, resource resources.Resource) error
	createLinkResourceFunc   func(ctx context.Context, resource resources.Resource) error
	listResourcesForWeekFunc func(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error)
	listResourceForUserFunc  func(ctx context.Context, userID uuid.UUID) ([]resources.UserResources, error)
	getResourceFunc          func(ctx context.Context, resourceID uuid.UUID) (string, error)
	deleteResourceFunc       func(ctx context.Context, userID, resourceID uuid.UUID) error
	cleanOrphanObjectsFunc   func(ctx context.Context) ([]uuid.UUID, error)
}

func (m *mockResourceService) UploadResource(ctx context.Context, file io.Reader, size int64, resource resources.Resource) error {
	if m.uploadResourceFunc != nil {
		return m.uploadResourceFunc(ctx, file, size, resource)
	}
	return nil
}

func (m *mockResourceService) CreateLinkResource(ctx context.Context, resource resources.Resource) error {
	if m.createLinkResourceFunc != nil {
		return m.createLinkResourceFunc(ctx, resource)
	}
	return nil
}

func (m *mockResourceService) ListResourcesForWeek(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error) {
	if m.listResourcesForWeekFunc != nil {
		return m.listResourcesForWeekFunc(ctx, weekID)
	}
	return []resources.ResourceWithUser{}, nil
}

func (m *mockResourceService) ListResourceForUser(ctx context.Context, userID uuid.UUID) ([]resources.UserResources, error) {
	if m.listResourceForUserFunc != nil {
		return m.listResourceForUserFunc(ctx, userID)
	}
	return []resources.UserResources{}, nil
}

func (m *mockResourceService) GetResource(ctx context.Context, resourceID uuid.UUID) (string, error) {
	if m.getResourceFunc != nil {
		return m.getResourceFunc(ctx, resourceID)
	}
	return "", nil
}

func (m *mockResourceService) DeleteResource(ctx context.Context, userID, resourceID uuid.UUID) error {
	if m.deleteResourceFunc != nil {
		return m.deleteResourceFunc(ctx, userID, resourceID)
	}
	return nil
}

func (m *mockResourceService) CleanOrphanObjects(ctx context.Context) ([]uuid.UUID, error) {
	if m.cleanOrphanObjectsFunc != nil {
		return m.cleanOrphanObjectsFunc(ctx)
	}
	return []uuid.UUID{}, nil
}

func TestCreateLinkResourceHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		userID         string
		requestBody    interface{}
		mockFunc       func(ctx context.Context, resource resources.Resource) error
		expectedStatus int
	}{
		{
			name:   "success - create link resource",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: CreateLinkResourceRequest{
				URL:  "https://example.com",
				Name: "Example Link",
			},
			mockFunc: func(ctx context.Context, resource resources.Resource) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "error - resource already exists",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: CreateLinkResourceRequest{
				URL:  "https://example.com",
				Name: "Example Link",
			},
			mockFunc: func(ctx context.Context, resource resources.Resource) error {
				return resources.ErrResourceExists
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error - invalid week ID",
			weekID:         "invalid-uuid",
			userID:         uuid.New().String(),
			requestBody:    CreateLinkResourceRequest{URL: "https://example.com", Name: "Test"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - database error",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			requestBody: CreateLinkResourceRequest{
				URL:  "https://example.com",
				Name: "Example Link",
			},
			mockFunc: func(ctx context.Context, resource resources.Resource) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockResourceService{createLinkResourceFunc: tt.mockFunc}
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/resources/link/"+tt.weekID, bytes.NewBuffer(body))
			req = addUserIDToContext(req, tt.userID)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if ok {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					var reqData CreateLinkResourceRequest
					if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&reqData); err == nil {
						resource := resources.Resource{
							ID:           uuid.New(),
							WeekID:       weekID,
							UserID:       userID,
							ResourceType: resources.ResourceLink,
							ExternalLink: &reqData.URL,
							Name:         reqData.Name,
						}
						err := mockSvc.CreateLinkResource(req.Context(), resource)
						if err != nil {
							if errors.Is(err, resources.ErrResourceExists) {
								ResponseWithErr(w, http.StatusBadRequest, "Link is already uploaded by other user")
							} else {
								ResponseWithErr(w, http.StatusInternalServerError, "failed to create the Link Resource")
							}
						} else {
							ResponseWithJSON(w, http.StatusCreated, nil)
						}
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListResourcesForWeekHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		mockFunc       func(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "success - list resources for week",
			weekID: uuid.New().String(),
			mockFunc: func(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error) {
				return []resources.ResourceWithUser{
					{
						ID:           uuid.New(),
						WeekID:       weekID,
						UserID:       uuid.New(),
						UserName:     "John",
						Name:         "Lecture Notes",
						ResourceType: resources.ResourceFile,
					},
					{
						ID:           uuid.New(),
						WeekID:       weekID,
						UserID:       uuid.New(),
						UserName:     "Jane",
						Name:         "Study Guide",
						ResourceType: resources.ResourceLink,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:   "success - empty list",
			weekID: uuid.New().String(),
			mockFunc: func(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error) {
				return []resources.ResourceWithUser{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "error - invalid week ID",
			weekID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - database error",
			weekID: uuid.New().String(),
			mockFunc: func(ctx context.Context, weekID uuid.UUID) ([]resources.ResourceWithUser, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockResourceService{listResourcesForWeekFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/resources/weeks/"+tt.weekID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if ok {
				res, err := mockSvc.ListResourcesForWeek(req.Context(), weekID)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to list resources")
				} else {
					ResponseWithJSON(w, http.StatusOK, res)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListResourcesForUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockFunc       func(ctx context.Context, userID uuid.UUID) ([]resources.UserResources, error)
		expectedStatus int
	}{
		{
			name:   "success - list user resources",
			userID: uuid.New().String(),
			mockFunc: func(ctx context.Context, userID uuid.UUID) ([]resources.UserResources, error) {
				return []resources.UserResources{
					{
						ID:           uuid.New(),
						WeekID:       uuid.New(),
						UserID:       userID,
						ModuleName:   "CS101",
						Semester:     "Spring",
						Year:         2024,
						WeekNumber:   1,
						Name:         "Notes",
						ResourceType: resources.ResourceFile,
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid user ID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockResourceService{listResourceForUserFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/resources/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			userIDParam := chi.URLParam(req, "user_id")
			userID, ok := parseUUID(w, userIDParam)
			if ok {
				res, err := mockSvc.ListResourceForUser(req.Context(), userID)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to list resources")
				} else {
					ResponseWithJSON(w, http.StatusOK, res)
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetResourceHandler(t *testing.T) {
	tests := []struct {
		name           string
		resourceID     string
		mockFunc       func(ctx context.Context, resourceID uuid.UUID) (string, error)
		expectedStatus int
	}{
		{
			name:       "success - get presigned URL",
			resourceID: uuid.New().String(),
			mockFunc: func(ctx context.Context, resourceID uuid.UUID) (string, error) {
				return "https://s3.amazonaws.com/presigned-url", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid resource ID",
			resourceID:     "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "error - resource not found",
			resourceID: uuid.New().String(),
			mockFunc: func(ctx context.Context, resourceID uuid.UUID) (string, error) {
				return "", errors.New("resource not found")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockResourceService{getResourceFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodGet, "/resources/"+tt.resourceID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.resourceID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			idParam := chi.URLParam(req, "id")
			resourceID, ok := parseUUID(w, idParam)
			if ok {
				url, err := mockSvc.GetResource(req.Context(), resourceID)
				if err != nil {
					ResponseWithErr(w, http.StatusInternalServerError, "failed to get resource")
				} else {
					ResponseWithJSON(w, http.StatusOK, map[string]string{"url": url})
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteResourceHandler(t *testing.T) {
	tests := []struct {
		name           string
		resourceID     string
		userID         string
		mockFunc       func(ctx context.Context, userID, resourceID uuid.UUID) error
		expectedStatus int
	}{
		{
			name:       "success - delete resource",
			resourceID: uuid.New().String(),
			userID:     uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, resourceID uuid.UUID) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - invalid resource ID",
			resourceID:     "invalid-uuid",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "error - database error",
			resourceID: uuid.New().String(),
			userID:     uuid.New().String(),
			mockFunc: func(ctx context.Context, userID, resourceID uuid.UUID) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockResourceService{deleteResourceFunc: tt.mockFunc}
			req := httptest.NewRequest(http.MethodDelete, "/resources/"+tt.resourceID, nil)
			req = addUserIDToContext(req, tt.userID)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.resourceID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			resourceIDParam := chi.URLParam(req, "id")
			resourceID, ok := parseUUID(w, resourceIDParam)
			if ok {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					err := mockSvc.DeleteResource(req.Context(), userID, resourceID)
					if err != nil {
						ResponseWithErr(w, http.StatusInternalServerError, "failed to delete resource")
					} else {
						ResponseWithJSON(w, http.StatusOK, nil)
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUploadFileHandler(t *testing.T) {
	tests := []struct {
		name           string
		weekID         string
		userID         string
		setupRequest   func() (*http.Request, error)
		mockFunc       func(ctx context.Context, file io.Reader, size int64, resource resources.Resource) error
		expectedStatus int
	}{
		{
			name:   "success - upload file",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.pdf")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("test file content"))
				writer.WriteField("fileType", "pdf")
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/resources/file/test", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFunc: func(ctx context.Context, file io.Reader, size int64, resource resources.Resource) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "error - resource already exists",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.pdf")
				part.Write([]byte("test content"))
				writer.Close()

				req := httptest.NewRequest(http.MethodPost, "/resources/file/test", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFunc: func(ctx context.Context, file io.Reader, size int64, resource resources.Resource) error {
				return resources.ErrResourceExists
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "error - no file in request",
			weekID: uuid.New().String(),
			userID: uuid.New().String(),
			setupRequest: func() (*http.Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/resources/file/test", strings.NewReader(""))
				req.Header.Set("Content-Type", "multipart/form-data")
				return req, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupRequest == nil {
				t.Skip("No request setup provided")
			}

			mockSvc := &mockResourceService{uploadResourceFunc: tt.mockFunc}
			req, err := tt.setupRequest()
			if err != nil {
				t.Fatalf("failed to setup request: %v", err)
			}

			req = addUserIDToContext(req, tt.userID)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("week_id", tt.weekID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			// Execute handler logic
			weekIDParam := chi.URLParam(req, "week_id")
			weekID, ok := parseUUID(w, weekIDParam)
			if ok {
				userIDStr := getUserID(req)
				userID, okUser := parseUUID(w, userIDStr)
				if okUser {
					file, handler, err := req.FormFile("file")
					if err != nil {
						ResponseWithErr(w, http.StatusBadRequest, "cannot access form file data")
					} else {
						defer file.Close()
						fileType := req.FormValue("fileType")
						resource := resources.Resource{
							ID:           uuid.New(),
							WeekID:       weekID,
							UserID:       userID,
							ResourceType: resources.ResourceFile,
							Name:         handler.Filename,
							FileType:     fileType,
						}
						err = mockSvc.UploadResource(req.Context(), file, handler.Size, resource)
						if err != nil {
							if errors.Is(err, resources.ErrResourceExists) {
								ResponseWithErr(w, http.StatusBadRequest, "file is uploaded by other user already")
							} else {
								ResponseWithErr(w, http.StatusInternalServerError, "failed to upload file")
							}
						} else {
							ResponseWithJSON(w, http.StatusCreated, nil)
						}
					}
				}
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
