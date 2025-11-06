package store

import (
	"errors"
	"library-api/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestBookStore_Create(t *testing.T) {
	testCases := []struct {
		description   string
		book          model.Book
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "create book successfully",
			book: model.Book{
				ID:        "0eabf8fc-1867-48c4-b835-271db2be1f2e",
				AuthorsID: "ce99cad9-9d1c-4e8c-a306-e51d7022926e",
				Title:     "Desert Stars",
				Genre:     "IT",
				ISBN:      "978-1-00001-000-1",
			},
			expectedError: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO books`).
					WithArgs("0eabf8fc-1867-48c4-b835-271db2be1f2e", "ce99cad9-9d1c-4e8c-a306-e51d7022926e", "Desert Stars", "IT", "978-1-00001-000-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description:   "error db",
			expectedError: errors.New("error"),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO books`).
					WillReturnError(errors.New("error"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewBookStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Create(&testCase.book)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNewBookStore_Get(t *testing.T) {
	columns := []string{"id", "authors_id", "title", "genre", "isbn"}
	testCases := []struct {
		description   string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  []model.Book
		expectedError error
	}{
		{
			description: "book store created successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow("0eabf8fc-1867-48c4-b835-271db2be1f2e", "ce99cad9-9d1c-4e8c-a306-e51d7022926e", "Desert Stars", "IT", "978-1-00001-000-1").
					AddRow("11f76f2b-9aa1-483c-91e4-3312b931e437", "ed6a7278-97a8-4382-847d-a4a0b02bca86", "Fictional Truths", "Fiction", "978-1-00002-000-1")

				mock.ExpectQuery("SELECT \\* FROM books").
					WillReturnRows(rows)
			},
			expectedBody: []model.Book{
				{
					ID:        "0eabf8fc-1867-48c4-b835-271db2be1f2e",
					AuthorsID: "ce99cad9-9d1c-4e8c-a306-e51d7022926e",
					Title:     "Desert Stars",
					Genre:     "IT",
					ISBN:      "978-1-00001-000-1",
				},
				{
					ID:        "11f76f2b-9aa1-483c-91e4-3312b931e437",
					AuthorsID: "ed6a7278-97a8-4382-847d-a4a0b02bca86",
					Title:     "Fictional Truths",
					Genre:     "Fiction",
					ISBN:      "978-1-00002-000-1",
				},
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM books").
					WillReturnError(errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		{
			description: "empty db",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"only one row"}).
					AddRow("hello")

				mock.ExpectQuery("SELECT \\* FROM books").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 1 destination arguments in Scan, not 5"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewBookStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			body, err := s.Get()
			assert.Equal(t, testCase.expectedError, err)

			assert.Equal(t, testCase.expectedBody, body)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBookStore_Exists(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "delete book successfully",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM books WHERE id = \\$1\\)").
					WithArgs("0eabf8fc-1867-48c4-b835-271db2be1f2e").
					WillReturnRows(rows)
			},
		},
		{
			description: "error db",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM books WHERE id = \\$1\\)").
					WillReturnError(errors.New("select error"))
			},
			expectedError: errors.New("select error"),
		},
		{
			description: "scan error",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1st row", "second row"}).
					AddRow("hello", "world")

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM books WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("sql: expected 2 destination arguments in Scan, not 1"),
		},
		{
			description: "book doesn't exist",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(false)

				mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM books WHERE id = \\$1\\)").
					WillReturnRows(rows)
			},
			expectedError: errors.New("book does not exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			s := NewBookStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Exists(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBookStore_Update(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		body          model.Book
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  model.Book
		expectedError error
	}{
		{
			description: "update book sucsessfully",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			body: model.Book{
				ID:        "0eabf8fc-1867-48c4-b835-271db2be1f2e",
				AuthorsID: "ce99cad9-9d1c-4e8c-a306-e51d7022926e",
				Title:     "Desert Stars",
				Genre:     "IT",
				ISBN:      "978-1-00001-000-1",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE books").
					WithArgs("ce99cad9-9d1c-4e8c-a306-e51d7022926e", "Desert Stars", "IT", "978-1-00001-000-1", "0eabf8fc-1867-48c4-b835-271db2be1f2e").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedBody: model.Book{
				ID:        "0eabf8fc-1867-48c4-b835-271db2be1f2e",
				AuthorsID: "ce99cad9-9d1c-4e8c-a306-e51d7022926e",
				Title:     "Desert Stars",
				Genre:     "IT",
				ISBN:      "978-1-00001-000-1",
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE books").
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

			s := NewBookStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Update(testCase.id, &testCase.body)
			assert.Equal(t, testCase.expectedError, err)

			assert.Equal(t, testCase.expectedBody, testCase.body)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBookStore_Delete(t *testing.T) {
	testCases := []struct {
		description   string
		id            string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "delete book successfully",
			id:          "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM books WHERE ID = \\$1").
					WithArgs("0eabf8fc-1867-48c4-b835-271db2be1f2e").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM books").
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

			s := NewBookStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Delete(testCase.id)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
