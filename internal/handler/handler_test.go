package handler

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorHandler(t *testing.T) {
	mockAuthorStore := new(MockAuthorStore)
	actualAuthorHandler := NewAuthorHandler(mockAuthorStore, hclog.NewNullLogger())

	expectedAuthorHandler := &AuthorHandler{
		store:  mockAuthorStore,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expectedAuthorHandler, actualAuthorHandler)
}

func TestNewBookHandler(t *testing.T) {
	mockBookStore := new(MockBookStore)
	actualBookHandler := NewBookHandler(mockBookStore, hclog.NewNullLogger())

	expectedBookHandler := &BookHandler{
		store:  mockBookStore,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expectedBookHandler, actualBookHandler)
}

func TestNewMemberHandler(t *testing.T) {
	mockMemberStore := new(MockMemberStore)
	actualMemberHandler := NewMemberHandler(mockMemberStore, hclog.NewNullLogger())

	expectedMemberHandler := &MemberHandler{
		store:  mockMemberStore,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expectedMemberHandler, actualMemberHandler)
}

func TestNewBorrowedHandler(t *testing.T) {
	mockBorrowedStore := new(MockBorrowedStore)
	actualBorrowedHandler := NewBorrowedHandler(mockBorrowedStore, hclog.NewNullLogger())

	expectedBorrowedHandler := &BorrowedHandler{
		store:  mockBorrowedStore,
		logger: hclog.NewNullLogger(),
	}

	assert.Equal(t, expectedBorrowedHandler, actualBorrowedHandler)
}
