package store

import (
	"library-api/internal/model"

	"github.com/lib/pq"
)

func (b *BorrowedStore) Create(book *model.Borrowed) error {
	_, err := b.db.Exec(`INSERT INTO borrowed_books (member_id, book_id) VALUES ($1, $2)`,
		&book.MemberID, &book.BookID)
	if err != nil {
		b.logger.Error("failed to create book",
			"member_id", book.MemberID,
			"book_id", book.BookID,
			"error", err.Error())
		return err
	}

	return nil
}

func (b *BorrowedStore) Get(id string) ([]model.Book, error) {
	rows, err := b.db.Query(`SELECT books.title, authors.full_name, books.genre, books.isbn
									FROM books, authors, borrowed_books
									WHERE (authors.id = books.authors_id 
									           AND books.id = borrowed_books.book_id 
									           AND borrowed_books.member_id = $1)`, id)
	if err != nil {
		b.logger.Error("get books failed for member",
			"id", id,
			"error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		var authorFullName *string

		err = rows.Scan(&book.Title, &authorFullName, &book.Genre, &book.ISBN)
		if err != nil {
			b.logger.Error("scanning selected failed for books of member",
				"id", id,
				"error", err.Error())
			return nil, err
		}

		book.Author = model.Author{
			FullName: authorFullName,
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *BorrowedStore) Delete(memberId string, bookId string) error {
	_, err := b.db.Exec(`DELETE FROM borrowed_books WHERE member_id = $1 AND book_id = $2`,
		memberId, bookId)
	if err != nil {
		b.logger.Error("delete book failed for member",
			"member_id", memberId,
			"book_id", bookId,
			"error", err.Error())
		return err
	}

	return nil
}

func (b *BorrowedStore) DeleteList(id string, books []string) error {
	_, err := b.db.Exec(`DELETE FROM borrowed_books WHERE (member_id = $1 AND book_id = ANY($2))`,
		id, pq.Array(books))
	if err != nil {
		b.logger.Error("delete list of books failed for member",
			"member_id", id,
			"error", err.Error())
		return err
	}
	return nil
}
