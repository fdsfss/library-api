package app

func (s *server) router() {
	s.app.Get("/healthz", s.healthcheck)

	s.app.Get("/authors", s.authorHandler.Get)
	s.app.Post("/author", s.authorHandler.Create)
	s.app.Patch("/author/:id", s.authorHandler.Update)
	s.app.Delete("/author/:id", s.authorHandler.Delete)
	s.app.Get("/author/:id/books", s.authorHandler.GetAuthorBooks)

	s.app.Get("/books", s.bookHandler.Get)
	s.app.Post("/book", s.bookHandler.Create)
	s.app.Patch("/book/:id", s.bookHandler.Update)
	s.app.Delete("/book/:id", s.bookHandler.Delete)

	s.app.Get("/members", s.memberHandler.Get)
	s.app.Post("/member", s.memberHandler.Create)
	s.app.Patch("/member/:id", s.memberHandler.Update)
	s.app.Delete("/member/:id", s.memberHandler.Delete)

	s.app.Get("/member/:id/borrowed", s.borrowedHandler.Get)
	s.app.Post("/member/borrowed", s.borrowedHandler.Create)
	s.app.Delete("/member/:id/borrowed/:book_id", s.borrowedHandler.Delete)
	s.app.Delete("/member/:id/borrowed", s.borrowedHandler.DeleteList)
}
