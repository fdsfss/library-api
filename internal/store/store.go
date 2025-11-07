package store

import (
	"database/sql"

	"github.com/hashicorp/go-hclog"
)

type AuthorStore struct {
	db     *sql.DB
	logger hclog.Logger
}

func NewAuthorStore(db *sql.DB, logger hclog.Logger) *AuthorStore {
	return &AuthorStore{
		db:     db,
		logger: logger,
	}
}

type BookStore struct {
	db     *sql.DB
	logger hclog.Logger
}

func NewBookStore(db *sql.DB, logger hclog.Logger) *BookStore {
	return &BookStore{
		db:     db,
		logger: logger,
	}
}

type BorrowedStore struct {
	db     *sql.DB
	logger hclog.Logger
}

func NewBorrowedStore(db *sql.DB, logger hclog.Logger) *BorrowedStore {
	return &BorrowedStore{
		db:     db,
		logger: logger,
	}
}

type MemberStore struct {
	db     *sql.DB
	logger hclog.Logger
}

func NewMemberStore(db *sql.DB, logger hclog.Logger) *MemberStore {
	return &MemberStore{
		db:     db,
		logger: logger,
	}
}
