package store

import (
	"errors"
	"library-api/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestBorrowedStore_Get(t *testing.T) {
	testCases := []struct {
		description   string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  []model.Member
		expectedError error
	}{
		{
			description: "get member successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "full_name"}).
					AddRow("3f45f596-ae05-4a60-802c-e2d45e7c26a2", "Samir Kenzhe").
					AddRow("8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4", "Amina Tulegen")

				mock.ExpectQuery("SELECT \\* FROM members").
					WillReturnRows(rows)
			},
			expectedBody: []model.Member{
				{
					ID:       "3f45f596-ae05-4a60-802c-e2d45e7c26a2",
					FullName: "Samir Kenzhe",
				},
				{
					ID:       "8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4",
					FullName: "Amina Tulegen",
				},
			},
		},
		{
			description: "db error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM members").
					WillReturnError(errors.New("select all failed for members"))
			},
			expectedError: errors.New("select all failed for members"),
		},
		{
			description: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1st row"}).
					AddRow("hello")

				mock.ExpectQuery("SELECT \\* FROM members").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 1 destination arguments in Scan, not 2"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewMemberStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			body, err := s.Get()
			assert.Equal(t, testCase.expectedBody, body)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNewMemberStore_Create(t *testing.T) {
	testCases := []struct {
		description   string
		body          model.Member
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "create member successfully",
			body: model.Member{
				ID:       "3f45f596-ae05-4a60-802c-e2d45e7c26a2",
				FullName: "Samir Kenzhe",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO members\\(id, full_name\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("3f45f596-ae05-4a60-802c-e2d45e7c26a2", "Samir Kenzhe").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			body: model.Member{
				ID:       "3f45f596-ae05-4a60-802c-e2d45e7c26a2",
				FullName: "Samir Kenzhe",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO members\\(id, full_name\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("3f45f596-ae05-4a60-802c-e2d45e7c26a2", "Samir Kenzhe").
					WillReturnError(errors.New("insert request failed"))
			},
			expectedError: errors.New("insert request failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewMemberStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Create(&testCase.body)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestMemberStore_Exists(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "member exists successfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM members WHERE id = \\$1\\)").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnRows(rows)
			},
		},
		{
			description: "error db",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM members WHERE id = \\$1\\)").
					WillReturnError(errors.New("select error"))
			},
			expectedError: errors.New("select error"),
		},
		{
			description: "scan error",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1st row", "second row"}).
					AddRow("hello", "world")

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM members WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 2 destination arguments in Scan, not 1"),
		},
		{
			description: "member doesn't exist",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(false)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM members WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("member does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewMemberStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Exists(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNewMemberStore_Update(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		body          model.Member
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "member updated successfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			body: model.Member{
				ID:       "b7eb3c06-6df8-4353-90f5-7ab897a77158",
				FullName: "John Doe",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE members SET full_name = \\$1 WHERE id = \\$2").
					WithArgs("John Doe", "b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "update request failed",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			body: model.Member{
				ID:       "b7eb3c06-6df8-4353-90f5-7ab897a77158",
				FullName: "John Doe",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE members SET full_name = \\$1 WHERE id = \\$2").
					WithArgs("John Doe", "b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnError(errors.New("update request failed"))
			},
			expectedError: errors.New("update request failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewMemberStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Update(testCase.id, &testCase.body)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestMemberStore_Delete(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "member deleted successfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM members WHERE id = \\$1").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "delete request failed",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM members WHERE id = \\$1").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnError(errors.New("delete request failed"))
			},
			expectedError: errors.New("delete request failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewMemberStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Delete(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
