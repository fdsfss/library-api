package store

import (
	"errors"
	"library-api/internal/model"
)

func (b *BookStore) Create(book *model.Book) error {
	_, err := b.db.Exec(`INSERT INTO books (id, authors_id, title, genre, isbn) 
								VALUES ($1, $2, $3, $4, $5)`,
		&book.ID, &book.AuthorsID, &book.Title, &book.Genre, &book.ISBN)
	if err != nil {
		b.logger.Error("failed to create book", "error", err.Error())
		return err
	}

	return nil
}

func (b *BookStore) Get() ([]model.Book, error) {
	rows, err := b.db.Query(`SELECT * FROM books`)
	if err != nil {
		b.logger.Error("select all failed for books", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		err = rows.Scan(&book.ID, &book.AuthorsID, &book.Title, &book.Genre, &book.ISBN)
		if err != nil {
			b.logger.Error("scanning selected failed for books", "error", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *BookStore) Exists(id string) error {
	rows, err := b.db.Query(`SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)`, id)
	if err != nil {
		b.logger.Info("id doesn't exists in books", "id", id, "info", err.Error())
		return err
	}

	for rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			b.logger.Info("scan error", "id", id, "error", err.Error())
			return err
		}

		if !exists {
			b.logger.Info("book does not exist", "id", id)
			return errors.New("book does not exist")
		}
	}

	return nil
}

func (b *BookStore) Update(id string, book *model.Book) error {
	_, err := b.db.Exec(`UPDATE books SET authors_id = $1, title = $2, genre = $3, ISBN = $4 WHERE id = $5`,
		&book.AuthorsID, &book.Title, &book.Genre, &book.ISBN, id)
	if err != nil {
		b.logger.Error("update failed for book", "id", id, "error", err.Error())
		return err
	}

	return nil
}

func (b *BookStore) Delete(id string) error {
	_, err := b.db.Exec(`DELETE FROM books WHERE ID = $1`, id)
	if err != nil {
		b.logger.Error("delete failed for books", "id", id, "error", err.Error())
		return err
	}
	return nil
}
