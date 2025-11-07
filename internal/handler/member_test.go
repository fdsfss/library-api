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

type MockMemberStore struct {
	mock.Mock
}

func (m *MockMemberStore) Create(member *model.Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberStore) Get() ([]model.Member, error) {
	args := m.Called()
	return args.Get(0).([]model.Member), args.Error(1)
}

func (m *MockMemberStore) Exists(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMemberStore) Update(id string, member *model.Member) error {
	args := m.Called(id, member)
	return args.Error(0)
}

func (m *MockMemberStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMemberHandler_Create(t *testing.T) {
	testCases := []struct {
		description    string
		body           any
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "member create success",
			body: model.Member{
				FullName: "John Doe",
			},
			expectedStatus: fiber.StatusCreated,
			expectedBody: fiber.Map{
				"message": "member created",
			},
		},
		{
			description:    "body parsing failed",
			body:           `{`,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "member creation failed",
			},
			expectedError: errors.New("body parsing fail"),
		},
		{
			description: "error from store",
			body: model.Member{
				FullName: "John Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "member creation failed",
			},
			expectedError: errors.New("member creation failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockMemberStore := new(MockMemberStore)
			memberHandler := &MemberHandler{
				store:  mockMemberStore,
				logger: hclog.NewNullLogger(),
			}

			app.Post("/member", memberHandler.Create)

			mockMemberStore.On("Create", mock.Anything).Return(testCase.expectedError)

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodPost, "/member", bytes.NewReader(body))
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

func TestMemberHandler_Get(t *testing.T) {
	testCases := []struct {
		description    string
		body           []model.Member
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description: "member get success",
			body: []model.Member{
				{
					ID:       "d4cc2192-9316-4f34-952f-9a0504f154d4",
					FullName: "John Doe",
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: []model.Member{
				{
					ID:       "d4cc2192-9316-4f34-952f-9a0504f154d4",
					FullName: "John Doe",
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
				"message": "no members found",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			var mockMemberStore MockMemberStore

			memberHandler := &MemberHandler{
				store:  &mockMemberStore,
				logger: hclog.NewNullLogger(),
			}

			app.Get("/members", memberHandler.Get)

			mockMemberStore.On("Get").Return(testCase.body, testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodGet, "/members", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if testCase.expectedStatus == fiber.StatusOK {
				var actual []model.Member
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

func TestMemberHandler_Update(t *testing.T) {
	testCases := []struct {
		description         string
		body                any
		expectedExistsError error
		expectedStatus      int
		expectedBody        any
		expectedUpdateError error
	}{
		{
			description: "member update success",
			body: model.Member{
				FullName: "John Doe",
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "member updated",
			},
		},
		{
			description:    "body parser error",
			body:           '{',
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"error": "member update failed",
			},
		},
		{
			description: "id doesn't exists",
			body: model.Member{
				FullName: "John Doe",
			},
			expectedExistsError: errors.New("member not found"),
			expectedStatus:      fiber.StatusNotFound,
			expectedBody: fiber.Map{
				"message": "member not found",
			},
		},
		{
			description: "error from store update",
			body: model.Member{
				FullName: "John Doe",
			},
			expectedStatus:      fiber.StatusBadRequest,
			expectedUpdateError: errors.New("member update failed"),
			expectedBody: fiber.Map{
				"error": "member update failed",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			app := fiber.New()

			mockMemberStore := new(MockMemberStore)
			memberHandler := &MemberHandler{
				store:  mockMemberStore,
				logger: hclog.NewNullLogger(),
			}

			app.Patch("/member/:id", memberHandler.Update)

			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			mockMemberStore.On("Exists", mock.Anything).Return(testCase.expectedExistsError).Once()

			mockMemberStore.On("Update", mock.Anything, mock.Anything).Return(testCase.expectedUpdateError).Once()

			req := httptest.NewRequest(fiber.MethodPatch, "/member/1de94d3e-09b2-4f62-bfff-964012c649d3", bytes.NewReader(body))
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

func TestMemberHandler_Delete(t *testing.T) {
	testCases := []struct {
		description    string
		expectedStatus int
		expectedBody   any
		expectedError  error
	}{
		{
			description:    "member delete success",
			expectedStatus: fiber.StatusOK,
			expectedBody: fiber.Map{
				"message": "member deleted",
			},
		},
		{
			description:    "member still has books, can't delete",
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: fiber.Map{
				"message": "member still has books, all books must be returned",
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

			mockMemberStore := new(MockMemberStore)
			memberHandler := &MemberHandler{
				store:  mockMemberStore,
				logger: hclog.NewNullLogger(),
			}

			app.Delete("/member/:id", memberHandler.Delete)

			mockMemberStore.On("Delete", mock.Anything).Return(testCase.expectedError).Once()

			req := httptest.NewRequest(fiber.MethodDelete, "/member/1de94d3e-09b2-4f62-bfff-964012c649d3", nil)
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
