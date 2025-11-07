package store

import (
	"errors"
	"library-api/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorStore_Get(t *testing.T) {
	fullName := "John Doe"
	columns := []string{"id", "full_name", "nick_name", "specialization"}
	testCases := []struct {
		description   string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  []model.Author
		expectedError error
	}{
		{
			description: "author store created successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow("b7eb3c06-6df8-4353-90f5-7ab897a77158", "John Doe", "johndoe123", "writer").
					AddRow("b44d8a61-6f6e-490e-88d6-45ff67088d0b", "John Doe", "johndoe123super", "writer")

				mock.ExpectQuery("SELECT \\* FROM authors").
					WillReturnRows(rows)
			},
			expectedBody: []model.Author{
				{
					ID:             "b7eb3c06-6df8-4353-90f5-7ab897a77158",
					FullName:       &fullName,
					NickName:       "johndoe123",
					Specialization: "writer",
				},
				{
					ID:             "b44d8a61-6f6e-490e-88d6-45ff67088d0b",
					FullName:       &fullName,
					NickName:       "johndoe123super",
					Specialization: "writer",
				},
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM authors").
					WillReturnError(errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			description: "empty db",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"only one row"}).
					AddRow("hello")

				mock.ExpectQuery("SELECT \\* FROM authors").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 1 destination arguments in Scan, not 4"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			body, err := s.Get()
			assert.Equal(t, testCase.expectedError, err)

			assert.Equal(t, testCase.expectedBody, body)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthorStore_Create(t *testing.T) {
	fullName := "John Doe"
	testCases := []struct {
		description   string
		author        model.Author
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "create author successfully",
			author: model.Author{
				ID:             "b7eb3c06-6df8-4353-90f5-7ab897a77158",
				FullName:       &fullName,
				NickName:       "johndoe123",
				Specialization: "writer",
			},
			expectedError: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO authors`).
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158", &fullName, "johndoe123", "writer").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description:   "error db",
			expectedError: errors.New("error"),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO authors`).
					WillReturnError(errors.New("error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Create(&testCase.author)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthorStore_Exists(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "author exists successfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM authors WHERE id = \\$1\\)").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnRows(rows)
			},
		},
		{
			description: "error db",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM authors WHERE id = \\$1\\)").
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

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM authors WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 2 destination arguments in Scan, not 1"),
		},
		{
			description: "author doesn't exist",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(false)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM authors WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("author does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Exists(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthorStore_Update(t *testing.T) {
	fullName := "John Doe"
	testCases := []struct {
		description   string
		id            string
		body          model.Author
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  model.Author
		expectedError error
	}{
		{
			description: "update author sucsessfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			body: model.Author{
				ID:             "b7eb3c06-6df8-4353-90f5-7ab897a77158",
				FullName:       &fullName,
				NickName:       "johndoe123 New",
				Specialization: "writer New",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE authors").
					WithArgs(&fullName, "johndoe123 New", "writer New", "b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedBody: model.Author{
				ID:             "b7eb3c06-6df8-4353-90f5-7ab897a77158",
				FullName:       &fullName,
				NickName:       "johndoe123 New",
				Specialization: "writer New",
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE authors").
					WillReturnError(errors.New("update failed"))
			},
			expectedError: errors.New("update failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Update(testCase.id, &testCase.body)
			assert.Equal(t, testCase.expectedError, err)

			assert.Equal(t, testCase.expectedBody, testCase.body)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthorStore_Delete(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "delete author successfully",
			id:          "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM author").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM author").
					WillReturnError(errors.New("delete failed"))
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Delete(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthorStore_GetAuthorBooks(t *testing.T) {
	testCases := []struct {
		description   string
		authorId      string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  []string
		expectedError error
	}{
		{
			description: "get author books successfully",
			authorId:    "b7eb3c06-6df8-4353-90f5-7ab897a77158",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"title"}).
					AddRow("title 1").
					AddRow("title 2")
				mock.ExpectQuery("SELECT title FROM books WHERE authors_id = \\$1").
					WithArgs("b7eb3c06-6df8-4353-90f5-7ab897a77158").
					WillReturnRows(rows)
			},
			expectedBody:  []string{"title 1", "title 2"},
			expectedError: nil,
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT title FROM books WHERE authors_id = \\$1").
					WillReturnError(errors.New("select for get for authors books failed"))
			},
			expectedError: errors.New("select for get for authors books failed"),
		},
		{
			description: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1st row", "second row"}).
					AddRow("hello", "world")
				mock.ExpectQuery("SELECT title FROM books WHERE authors_id = \\$1").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 2 destination arguments in Scan, not 1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewAuthorStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			res, err := s.GetAuthorsBooks(testCase.authorId)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedBody, res)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
