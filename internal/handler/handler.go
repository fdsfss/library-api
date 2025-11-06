package handler

import "github.com/hashicorp/go-hclog"

type AuthorHandler struct {
	store  authorStore
	logger hclog.Logger
}

func NewAuthorHandler(store authorStore, logger hclog.Logger) *AuthorHandler {
	return &AuthorHandler{
		store:  store,
		logger: logger,
	}
}

type BookHandler struct {
	store  bookStore
	logger hclog.Logger
}

func NewBookHandler(store bookStore, logger hclog.Logger) *BookHandler {
	return &BookHandler{
		store:  store,
		logger: logger,
	}
}

type BorrowedHandler struct {
	store  borrowedStore
	logger hclog.Logger
}

func NewBorrowedHandler(store borrowedStore, logger hclog.Logger) *BorrowedHandler {
	return &BorrowedHandler{
		store:  store,
		logger: logger,
	}
}

type MemberHandler struct {
	store  memberStore
	logger hclog.Logger
}

func NewMemberHandler(store memberStore, logger hclog.Logger) *MemberHandler {
	return &MemberHandler{
		store:  store,
		logger: logger,
	}
}
