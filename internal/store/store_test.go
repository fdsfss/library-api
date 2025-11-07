package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestAuthorStore(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDb.Close()

	actual := NewAuthorStore(mockDb, hclog.NewNullLogger())

	expected := &AuthorStore{
		db:     mockDb,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expected, actual)
}

func TestBookStore(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDb.Close()

	actual := NewBookStore(mockDb, hclog.NewNullLogger())

	expected := &BookStore{
		db:     mockDb,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expected, actual)
}

func TestBorrowedStore(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDb.Close()

	actual := NewBorrowedStore(mockDb, hclog.NewNullLogger())

	expected := &BorrowedStore{
		db:     mockDb,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberStore(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDb.Close()

	actual := NewMemberStore(mockDb, hclog.NewNullLogger())

	expected := &MemberStore{
		db:     mockDb,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expected, actual)
}
