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

type MockAuthorStore struct {
	mock.Mock
}

func (m *MockAuthorStore) Create(author *model.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorStore) Get() ([]model.Author, error) {
	args := m.Called()
	return args.Get(0).([]model.Author), args.Error(1)
}

func (m *MockAuthorStore) Exists(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAuthorStore) Update(id string, author *model.Author) error {
	args := m.Called(id, author)
	return args.Error(0)
}

func (m *MockAuthorStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAuthorStore) GetAuthorsBooks(id string) ([]string, error) {
	args := m.Called(id)
	return args.Get(0).([]string), args.Error(1)
}

func TestAuthorHandler_Create(t *testing.T) {
	fullName := "John Doe"
	testCases := []struct {
		description    string
		body           any
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "author successfully created",
			body: model.Author{
				FullName:       &fullName,
				NickName:       "johndoe123",
				Specialization: "Writer",
			},
			expectedStatus: fiber.StatusCreated,
			expectedBody: fiber.Map{
				"message": "author created",
			},
		},
		{
			description:    "body parsing failed",
			body:           `{`,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "author creation failed",
			},
			expectedError: errors.New("body parsing fail"),
		},
		{
			description: "error from store",
			body: model.Author{
				FullName: &fullName,
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "author creation failed",
			},
			expectedError: errors.New("author creation failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockAuthorStore := new(MockAuthorStore)
			authorHandler := &AuthorHandler{
				store:  mockAuthorStore,
				logger: hclog.NewNullLogger(),
			}

			app.Post("/author", authorHandler.Create)

			mockAuthorStore.On("Create", mock.Anything).Return(testCase.expectedError).Once()

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/author", bytes.NewReader(body))
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

func TestAuthorHandler_Get(t *testing.T) {
	fullName := "John Doe"
	testCases := []struct {
		description    string
		body           []model.Author
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "author get success",
			body: []model.Author{
				{
					ID:             "6da2643a-a24a-4b85-a188-2dd5502eaa66",
					FullName:       &fullName,
					NickName:       "johndoe123",
					Specialization: "Writer",
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: []model.Author{
				{
					ID:             "6da2643a-a24a-4b85-a188-2dd5502eaa66",
					FullName:       &fullName,
					NickName:       "johndoe123",
					Specialization: "Writer",
				},
			},
		},
		{
			description:    "store error",
			body:           nil,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: fiber.Map{
				"error": "server error",
			},
			expectedError: errors.New("server error"),
		},
		{
			description:    "empty db",
			body:           []model.Author{},
			expectedStatus: fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"message": "authors not found",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			var mockAuthorStore MockAuthorStore

			authorHandler := &AuthorHandler{
				store:  &mockAuthorStore,
				logger: hclog.NewNullLogger(),
			}

			app.Get("/authors", authorHandler.Get)

			mockAuthorStore.On("Get").Return(testCase.body, testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodGet, "/authors", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if resp.StatusCode == fiber.StatusOK {
				var actual []model.Author
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

func TestAuthorHandler_Update(t *testing.T) {
	fullName := "John Doe"
	testCases := []struct {
		description         string
		body                any
		expectedExistsError error
		expectedStatus      int
		expectedBody        any
		expectedUpdateError error
	}{
		{
			description: "author update success",
			body: model.Author{
				FullName:       &fullName,
				NickName:       "johndoe123",
				Specialization: "Writer",
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "author updated",
			},
		},
		{
			description:    "body parser error",
			body:           '{',
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "author update failed",
			},
		},
		{
			description: "id doesn't exists",
			body: model.Author{
				FullName:       &fullName,
				NickName:       "johndoe123",
				Specialization: "Writer",
			},
			expectedExistsError: errors.New("author not found"),
			expectedStatus:      fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"error": "author not found",
			},
		},
		{
			description: "error from store",
			body: model.Author{
				FullName: &fullName,
				NickName: "johndoe123",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "author update failed",
			},
			expectedUpdateError: errors.New("author update failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockAuthorStore := new(MockAuthorStore)
			authorHandler := &AuthorHandler{
				store:  mockAuthorStore,
				logger: hclog.NewNullLogger(),
			}

			app.Patch("/author/:id", authorHandler.Update)

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			mockAuthorStore.On("Exists", mock.Anything).Return(testCase.expectedExistsError).Once()

			mockAuthorStore.On("Update", mock.Anything, mock.Anything).Return(testCase.expectedUpdateError).Once()

			req := httptest.NewRequest(fiber.MethodPatch, "/author/4dbec5df-c354-4c0a-8f33-7832dfbc12c0", bytes.NewReader(body))
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

func TestAuthorHandler_Delete(t *testing.T) {
	testCases := []struct {
		description    string
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description:    "author delete success",
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "author deleted",
			},
		},
		{
			description:    "author still has books, can't delete",
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"message": "author has related recordings and cannot be deleted",
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

			mockAuthorStore := new(MockAuthorStore)
			authorHandler := &AuthorHandler{
				store:  mockAuthorStore,
				logger: hclog.NewNullLogger(),
			}

			app.Delete("/author/:id", authorHandler.Delete)

			mockAuthorStore.On("Delete", mock.Anything).Return(testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodDelete, "/author/4dbec5df-c354-4c0a-8f33-7832dfbc12c0", nil)
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

func TestAuthorHandler_GetAuthorBooks(t *testing.T) {
	testCases := []struct {
		description    string
		body           []string
		expectedError  error
		expectedBody   any
		expectedStatus int
	}{
		{
			description:    "author get books success",
			body:           []string{"book1", "book2"},
			expectedBody:   []string{"book1", "book2"},
			expectedStatus: fiber.StatusOK,
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
			description:    "author has no books",
			expectedStatus: fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"error": "book not found",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockAuthorStore := new(MockAuthorStore)
			authorHandler := &AuthorHandler{
				store:  mockAuthorStore,
				logger: hclog.NewNullLogger(),
			}

			app.Get("/author/:id/books", authorHandler.GetAuthorBooks)

			mockAuthorStore.On("GetAuthorsBooks", mock.Anything).Return(testCase.body, testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodGet, "/author/4dbec5df-c354-4c0a-8f33-7832dfbc12c0/books", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if testCase.expectedStatus == fiber.StatusOK {
				var actual []string
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
