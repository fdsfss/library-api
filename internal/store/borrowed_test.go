package store

import (
	"errors"
	"library-api/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-hclog"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestBorrowedStore_Create(t *testing.T) {
	testCases := []struct {
		description   string
		body          model.Borrowed
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "borrowed store created successfully",
			body: model.Borrowed{
				MemberID: "dd2346fc-51c3-420f-a37e-8273d65120ad",
				BookID:   "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO borrowed_books \\(member_id, book_id\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", "0eabf8fc-1867-48c4-b835-271db2be1f2e").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			body: model.Borrowed{
				MemberID: "dd2346fc-51c3-420f-a37e-8273d65120ad",
				BookID:   "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO borrowed_books \\(member_id, book_id\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", "0eabf8fc-1867-48c4-b835-271db2be1f2e").
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

			s := NewBorrowedStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Create(&testCase.body)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNewBorrowedStore_Get(t *testing.T) {
	columns := []string{"title", "genre", "isbn", "authors_full_name"}
	authorsFullName := "Alice Johnson"
	testCases := []struct {
		description   string
		memberId      string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedBody  []model.Book
		expectedError error
	}{
		{
			description: "borrowed store created successfully",
			memberId:    "8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow("Fictional Truths", "Alice Johnson", "Fiction", "978-1-00002-000-1").
					AddRow("Ink and Imagination", "Alice Johnson", "Fiction", "978-1-00002-000-3")

				mock.ExpectQuery("SELECT books.title, authors.full_name, books.genre, books.isbn " +
					"FROM books, authors, borrowed_books " +
					"WHERE \\(authors.id = books.authors_id " +
					"AND books.id = borrowed_books.book_id " +
					"AND borrowed_books.member_id = \\$1\\)").
					WithArgs("8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4").
					WillReturnRows(rows)
			},
			expectedBody: []model.Book{
				{
					Title: "Fictional Truths",
					Genre: "Fiction",
					ISBN:  "978-1-00002-000-1",
					Author: model.Author{
						FullName: &authorsFullName,
					},
				},
				{
					Title: "Ink and Imagination",
					Genre: "Fiction",
					ISBN:  "978-1-00002-000-3",
					Author: model.Author{
						FullName: &authorsFullName,
					},
				},
			},
		},
		{
			description: "select error",
			memberId:    "8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT books.title, authors.full_name, books.genre, books.isbn " +
					"FROM books, authors, borrowed_books " +
					"WHERE \\(authors.id = books.authors_id " +
					"AND books.id = borrowed_books.book_id " +
					"AND borrowed_books.member_id = \\$1\\)").
					WithArgs("8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4").
					WillReturnError(errors.New("select error"))
			},
			expectedError: errors.New("select error"),
		},
		{
			description: "scan rows error",
			memberId:    "8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1st column"}).
					AddRow("Fictional Truths")

				mock.ExpectQuery("SELECT books.title, authors.full_name, books.genre, books.isbn " +
					"FROM books, authors, borrowed_books " +
					"WHERE \\(authors.id = books.authors_id " +
					"AND books.id = borrowed_books.book_id " +
					"AND borrowed_books.member_id = \\$1\\)").
					WithArgs("8ed3d7fd-88e6-44d9-b34b-9257a9a2d5b4").
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

			s := NewBorrowedStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			body, err := s.Get(testCase.memberId)
			assert.Equal(t, testCase.expectedError, err)

			assert.Equal(t, testCase.expectedBody, body)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBorrowedStore_Delete(t *testing.T) {
	testCases := []struct {
		description   string
		memberId      string
		bookId        string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "borrowed book deleted successfully",
			memberId:    "dd2346fc-51c3-420f-a37e-8273d65120ad",
			bookId:      "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM borrowed_books").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", "0eabf8fc-1867-48c4-b835-271db2be1f2e").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			memberId:    "dd2346fc-51c3-420f-a37e-8273d65120ad",
			bookId:      "0eabf8fc-1867-48c4-b835-271db2be1f2e",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM borrowed_books").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", "0eabf8fc-1867-48c4-b835-271db2be1f2e").
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

			s := NewBorrowedStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.Delete(testCase.memberId, testCase.bookId)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBorrowedStore_DeleteList(t *testing.T) {
	books := []string{"0eabf8fc-1867-48c4-b835-271db2be1f2e", "81790db6-a440-48e2-9951-d5fcf359fd7c"}

	testCases := []struct {
		description   string
		memberId      string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			description: "borrowed book deleted successfully",
			memberId:    "dd2346fc-51c3-420f-a37e-8273d65120ad",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM borrowed_books WHERE \\(member_id = \\$1 AND book_id = ANY\\(\\$2\\)\\)").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", pq.Array(books)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			description: "error db",
			memberId:    "dd2346fc-51c3-420f-a37e-8273d65120ad",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM borrowed_books WHERE \\(member_id = \\$1 AND book_id = ANY\\(\\$2\\)\\)").
					WithArgs("dd2346fc-51c3-420f-a37e-8273d65120ad", pq.Array(books)).
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

			s := NewBorrowedStore(db, hclog.NewNullLogger())

			testCase.setupMock(mock)

			err = s.DeleteList(testCase.memberId, books)
			assert.Equal(t, testCase.expectedError, err)

			err = mock.ExpectationsWereMet()
		})
	}
}
