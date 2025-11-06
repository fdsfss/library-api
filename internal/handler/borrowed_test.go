package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"library-api/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBorrowedStore struct {
	mock.Mock
}

func (m *MockBorrowedStore) Create(borrowed *model.Borrowed) error {
	args := m.Called(borrowed)
	return args.Error(0)
}

func (m *MockBorrowedStore) Get(id string) ([]model.Book, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBorrowedStore) Delete(memberId string, bookId string) error {
	args := m.Called(memberId, bookId)
	return args.Error(0)
}

func (m *MockBorrowedStore) DeleteList(memberId string, books []string) error {
	args := m.Called(memberId, books)
	return args.Error(0)
}

func TestBorrowedHandler_Create(t *testing.T) {
	testCases := []struct {
		description    string
		body           any
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "borrowed successfully created",
			body: model.Borrowed{
				MemberID: "3c864c77-39a5-4157-9fb6-39d72be81669",
				BookID:   "5dee5c81-5ee4-44a9-97e5-0eb7955792a4",
			},
			expectedStatus: fiber.StatusCreated,
			expectedBody: fiber.Map{
				"message": "borrowed book created",
			},
		},
		{
			description:    "body parsing failed",
			body:           `{`,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "borrowed book creation failed",
			},
			expectedError: errors.New("body parsing failed"),
		},
		{
			description: "error from store",
			body: model.Borrowed{
				MemberID: "3c864c77-39a5-4157-9fb6-39d72be81669",
				BookID:   "5dee5c81-5ee4-44a9-97e5-0eb7955792a4",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "borrowed book creation failed",
			},
			expectedError: errors.New("borrowed creation failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockBorrowedStore := new(MockBorrowedStore)
			borrowedHandler := &BorrowedHandler{
				store:  mockBorrowedStore,
				logger: hclog.NewNullLogger(),
			}

			app.Post("/member/borrowed", borrowedHandler.Create)

			mockBorrowedStore.On("Create", mock.Anything).Return(testCase.expectedError).Once()

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/member/borrowed", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var actual fiber.Map
			err = json.Unmarshal(respBody, &actual)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, actual)
		})
	}
}

func TestBorrowedHandler_Get(t *testing.T) {
	testCases := []struct {
		description    string
		body           []model.Book
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "borrowed get success",
			body: []model.Book{
				{
					ID:        "2d286219-8d2a-4b46-bec5-338d0ae1599a",
					AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
					Title:     "perfect book title",
					Genre:     "fantasy",
					ISBN:      "978-3-16-148410-0",
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: []model.Book{
				{
					ID:        "2d286219-8d2a-4b46-bec5-338d0ae1599a",
					AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
					Title:     "perfect book title",
					Genre:     "fantasy",
					ISBN:      "978-3-16-148410-0",
				},
			},
		},
		{
			description:    "store error",
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: fiber.Map{
				"error": "server error",
			},
			expectedError: errors.New("server error"),
		},
		{
			description:    "empty db",
			expectedStatus: fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"message": "no books found for this member",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			var mockBorrowedStore MockBorrowedStore

			borrowedHandler := &BorrowedHandler{
				store:  &mockBorrowedStore,
				logger: hclog.NewNullLogger(),
			}

			app.Get("/member/:id/borrowed", borrowedHandler.Get)

			mockBorrowedStore.On("Get", mock.Anything).Return(testCase.body, testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodGet, "/member/1de94d3e-09b2-4f62-bfff-964012c649d3/borrowed", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if testCase.expectedStatus == fiber.StatusOK {
				var actual []model.Book
				err = json.Unmarshal(respBody, &actual)
				assert.NoError(t, err)

				assert.Equal(t, testCase.expectedBody, actual)
			} else {
				var actual fiber.Map
				err = json.Unmarshal(respBody, &actual)
				assert.NoError(t, err)

				assert.Equal(t, testCase.expectedBody, actual)
			}
		})
	}
}

func TestBorrowedHandler_Delete(t *testing.T) {
	testCases := []struct {
		description    string
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description:    "borrowed delete success",
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "borrowed book deleted",
			},
		},
		{
			description:    "store error",
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: fiber.Map{
				"error": "server error",
			},
			expectedError: errors.New("server error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockBorrowedStore := new(MockBorrowedStore)
			borrowedHandler := &BorrowedHandler{
				store:  mockBorrowedStore,
				logger: hclog.NewNullLogger(),
			}

			app.Delete("/member/:id/borrowed/:book_id", borrowedHandler.Delete)

			mockBorrowedStore.On("Delete", mock.Anything, mock.Anything).Return(testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodDelete, "/member/1de94d3e-09b2-4f62-bfff-964012c649d3/borrowed/90a5d5a9-1161-4529-9841-0adb40a9eff1", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var actual fiber.Map
			err = json.Unmarshal(respBody, &actual)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, actual)
		})
	}
}

func TestBorrowedHandler_DeleteList(t *testing.T) {
	testCases := []struct {
		description    string
		body           any
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description:    "borrowed delete list success",
			body:           []string{"1652979f-bc46-44ea-ba2b-51e08f608021", "81790db6-a440-48e2-9951-d5fcf359fd7c"},
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "borrowed books deleted",
			},
		},
		{
			description:    "body parser error",
			body:           '{',
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "borrowed book delete failed",
			},
		},
		{
			description:    "store error",
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: fiber.Map{
				"error": "server error",
			},
			expectedError: errors.New("server error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockBorrowedStore := new(MockBorrowedStore)
			borrowedHandler := &BorrowedHandler{
				store:  mockBorrowedStore,
				logger: hclog.NewNullLogger(),
			}

			app.Delete("/member/:id/borrowed", borrowedHandler.DeleteList)

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			mockBorrowedStore.On("DeleteList", mock.Anything, mock.Anything).Return(testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodDelete, "/member/1de94d3e-09b2-4f62-bfff-964012c649d3/borrowed", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var actual fiber.Map
			err = json.Unmarshal(respBody, &actual)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedBody, actual)
		})
	}
}
