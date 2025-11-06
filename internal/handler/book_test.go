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

type MockBookStore struct {
	mock.Mock
}

func (m *MockBookStore) Create(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookStore) Get() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookStore) Exists(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookStore) Update(id string, book *model.Book) error {
	args := m.Called(id, book)
	return args.Error(0)
}

func (m *MockBookStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookHandler_Create(t *testing.T) {
	testCases := []struct {
		description    string
		body           any
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "book successfully created",
			body: model.Book{
				AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
				Title:     "perfect book title",
				Genre:     "fantasy",
				ISBN:      "978-3-16-148410-0",
			},
			expectedStatus: fiber.StatusCreated,
			expectedBody: fiber.Map{
				"id":      "dynamic ...",
				"message": "book created",
			},
		},
		{
			description:    "body parsing failed",
			body:           `{`,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "book creation failed",
			},
			expectedError: errors.New("body parsing fail"),
		},
		{
			description: "error from store",
			body: model.Book{
				Title: "perfect book title",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "book creation failed",
			},
			expectedError: errors.New("book creation failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockBookStore := new(MockBookStore)
			bookHandler := &BookHandler{
				store:  mockBookStore,
				logger: hclog.NewNullLogger(),
			}

			app.Post("/book", bookHandler.Create)

			mockBookStore.On("Create", mock.Anything).Return(testCase.expectedError).Once()

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/book", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var actual fiber.Map
			err = json.Unmarshal(respBody, &actual)
			assert.NoError(t, err)

			if testCase.expectedStatus != fiber.StatusCreated {
				assert.Equal(t, testCase.expectedBody, actual)
			} else {
				assert.Equal(t, testCase.expectedBody.(fiber.Map)["message"].(string), actual["message"].(string))
			}
		})
	}
}

func TestBookHandler_Get(t *testing.T) {
	testCases := []struct {
		description    string
		body           []model.Book
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "book get success",
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
				"message": "no books found",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			var mockBookStore MockBookStore

			bookHandler := &BookHandler{
				store:  &mockBookStore,
				logger: hclog.NewNullLogger(),
			}

			app.Get("/books", bookHandler.Get)

			mockBookStore.On("Get").Return(testCase.body, testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodGet, "/books", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if testCase.expectedStatus == fiber.StatusOK { //
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

func TestBookHandler_Update(t *testing.T) {
	testCases := []struct {
		description         string
		id                  string
		body                any
		expectedExistsError error
		expectedStatus      int
		expectedBody        any
		expectedUpdateError error
	}{
		{
			description: "book update success",
			body: model.Book{
				ID:        "2d286219-8d2a-4b46-bec5-338d0ae1599a",
				AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
				Title:     "perfect book title",
				Genre:     "fantasy",
				ISBN:      "978-3-16-148410-0",
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "book updated",
			},
		},
		{
			description:    "body parser error",
			body:           '{',
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "book update failed",
			},
		},
		{
			description: "id doesn't exists",
			body: model.Book{
				AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
				Title:     "perfect book title",
				Genre:     "fantasy",
				ISBN:      "978-3-16-148410-0",
			},
			expectedExistsError: errors.New("book not found"),
			expectedStatus:      fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"error": "book not found",
			},
		},
		{
			description: "error from store update",
			body: model.Book{
				AuthorsID: "c3690e20-5950-4a41-aa68-13f0791cdf98",
				Title:     "perfect book title",
				Genre:     "fantasy",
				ISBN:      "978-3-16-148410-0",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "book update failed",
			},
			expectedUpdateError: errors.New("book update failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockBookStore := new(MockBookStore)
			bookHandler := &BookHandler{
				store:  mockBookStore,
				logger: hclog.NewNullLogger(),
			}

			app.Patch("/book/:id", bookHandler.Update)

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			mockBookStore.On("Exists", mock.Anything).Return(testCase.expectedExistsError).Once()

			mockBookStore.On("Update", mock.Anything, mock.Anything).Return(testCase.expectedUpdateError).Once()

			req := httptest.NewRequest(fiber.MethodPatch, "/book/235fcd0e-98af-4af5-b985-68dab66085e1", bytes.NewReader(body))
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

func TestBookHandler_Delete(t *testing.T) {
	testCases := []struct {
		description    string
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description:    "book delete success",
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "book deleted",
			},
		},
		{
			description:    "book still has books, can't delete",
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"message": "book has related recordings and cannot be deleted",
			},
			expectedError: errors.New("violates foreign key constraint"),
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

			mockBookStore := new(MockBookStore)
			bookHandler := &BookHandler{
				store:  mockBookStore,
				logger: hclog.NewNullLogger(),
			}

			app.Delete("/book/:id", bookHandler.Delete)

			mockBookStore.On("Delete", mock.Anything).Return(testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodDelete, "/book/235fcd0e-98af-4af5-b985-68dab66085e1", nil)
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
