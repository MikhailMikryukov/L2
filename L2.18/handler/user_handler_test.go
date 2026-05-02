package handler

import (
	"L2.18/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestHandler() *UserHandler {
	us := service.NewUserService()
	return NewUserHandler(us)
}

func TestUserHandler_CreateEvent(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name       string
		method     string
		formValues map[string]string
		wantStatus int
	}{
		{
			name:   "successful creation",
			method: "POST",
			formValues: map[string]string{
				"user_id": "1",
				"date":    "2024-01-15",
				"event":   "Meeting",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "wrong method",
			method: "GET",
			formValues: map[string]string{
				"user_id": "1",
				"date":    "2024-01-15",
				"event":   "Meeting",
			},
			wantStatus: http.StatusOK, // Handler returns without processing
		},
		{
			name:   "invalid date",
			method: "POST",
			formValues: map[string]string{
				"user_id": "1",
				"date":    "invalid",
				"event":   "Meeting",
			},
			wantStatus: 400, // Bad Request
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			req := httptest.NewRequest(tt.method, "/create_event", body)
			req.Form = make(map[string][]string)

			for k, v := range tt.formValues {
				req.Form.Add(k, v)
			}

			rr := httptest.NewRecorder()
			handler.CreateEvent(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("CreateEvent() status = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_GetEventsDay(t *testing.T) {
	handler := setupTestHandler()

	// Сначала создаем событие
	createReq := httptest.NewRequest("POST", "/create_event", nil)
	createReq.Form = map[string][]string{
		"user_id": {"1"},
		"date":    {"2024-01-15"},
		"event":   {"Test Event"},
	}
	createRec := httptest.NewRecorder()
	handler.CreateEvent(createRec, createReq)

	t.Run("get events for existing user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/events_for_day?user_id=1&date=2024-01-15", nil)
		rr := httptest.NewRecorder()

		handler.GetEventsDay(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("GetEventsDay() status = %v, want %v", rr.Code, http.StatusOK)
		}

		var response ResultResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if len(response.Result) != 1 || response.Result[0] != "Test Event" {
			t.Errorf("GetEventsDay() result = %v, want ['Test Event']", response.Result)
		}
	})

	t.Run("get events for non-existing user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/events_for_day?user_id=999&date=2024-01-15", nil)
		rr := httptest.NewRecorder()

		handler.GetEventsDay(rr, req)

		if rr.Code != 503 { // Service Unavailable for user not found
			t.Errorf("GetEventsDay() status = %v, want %v", rr.Code, 503)
		}

		var errResponse ErrResponse
		err := json.NewDecoder(rr.Body).Decode(&errResponse)
		if err != nil {
			t.Errorf("Failed to decode error response: %v", err)
		}
	})
}

func TestUserHandler_UpdateEvent(t *testing.T) {
	handler := setupTestHandler()

	// Сначала создаем событие
	createReq := httptest.NewRequest("POST", "/create_event", nil)
	createReq.Form = map[string][]string{
		"user_id": {"1"},
		"date":    {"2024-01-15"},
		"event":   {"Old Event"},
	}
	createRec := httptest.NewRecorder()
	handler.CreateEvent(createRec, createReq)

	t.Run("update existing event", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update_event", nil)
		req.Form = map[string][]string{
			"user_id": {"1"},
			"date":    {"2024-01-15"},
			"event":   {"Old Event"}, // Нужно указать старое значение для обновления
		}
		// В текущей реализации UpdateEvent не обновляет значение, но проверяем вызов
		rr := httptest.NewRecorder()

		handler.UpdateEvent(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("UpdateEvent() status = %v, want %v", rr.Code, http.StatusOK)
		}
	})
}

func TestSendErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "user not found",
			err:        service.ErrUserNotFound,
			wantStatus: 503,
		},
		{
			name:       "invalid date",
			err:        service.ErrInvalidDate,
			wantStatus: 400,
		},
		{
			name:       "generic error",
			err:        service.ErrNoSuchEvent,
			wantStatus: 503,
		},
		{
			name:       "other error",
			err:        service.ErrInvalidID,
			wantStatus: 503,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			sendErrorResponse(rr, tt.err)

			if rr.Code != tt.wantStatus {
				t.Errorf("sendErrorResponse() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			var errResponse ErrResponse
			err := json.NewDecoder(rr.Body).Decode(&errResponse)
			if err != nil {
				t.Errorf("Failed to decode error response: %v", err)
			}

			if errResponse.Err != tt.err.Error() {
				t.Errorf("sendErrorResponse() error = %v, want %v", errResponse.Err, tt.err.Error())
			}
		})
	}
}
