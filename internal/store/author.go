package store

import (
	"errors"
	"library-api/internal/model"
)

func (a *AuthorStore) Get() ([]model.Author, error) {
	rows, err := a.db.Query(`SELECT * FROM authors`)
	if err != nil {
		a.logger.Error("select all request failed for author", "error", err.Error()) //
		return nil, err
	}
	defer rows.Close()

	var authors []model.Author
	for rows.Next() {
		var author model.Author
		err = rows.Scan(&author.ID, &author.FullName, &author.NickName, &author.Specialization)
		if err != nil {
			a.logger.Error("scanning selected failed for authors", "error", err.Error())
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (a *AuthorStore) Create(author *model.Author) error {
	_, err := a.db.Exec(`INSERT INTO authors (id, full_name, nick_name, specialization) VALUES ($1, $2, $3, $4)`,
		&author.ID, &author.FullName, &author.NickName, &author.Specialization)
	if err != nil {
		a.logger.Error("failed to create author", "error", err.Error())
		return err
	}

	return nil
}

func (a *AuthorStore) Exists(id string) error {
	rows, err := a.db.Query(`SELECT EXISTS (SELECT 1 FROM authors WHERE id = $1)`, id)
	if err != nil {
		a.logger.Info("select error", "id", id, "error", err.Error())
		return err
	}

	for rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			a.logger.Info("scan error", "id", id, "error", err.Error())
			return err
		}

		if !exists {
			a.logger.Info("author does not exist", "id", id)
			return errors.New("author does not exist")
		}
	}

	return nil
}

func (a *AuthorStore) Update(id string, author *model.Author) error {
	_, err := a.db.Exec(`UPDATE authors SET full_name = $1, nick_name = $2, specialization = $3 WHERE ID = $4`,
		&author.FullName, &author.NickName, &author.Specialization, id)
	if err != nil {
		a.logger.Error("update failed for author", "id", id, "error", err.Error())
		return err
	}

	return nil
}

func (a *AuthorStore) Delete(id string) error {
	_, err := a.db.Exec(`DELETE FROM authors WHERE ID = $1`, id)
	if err != nil {
		a.logger.Error("delete failed for authors", "id", id, "error", err.Error())
		return err
	}

	return nil
}

func (a *AuthorStore) GetAuthorsBooks(id string) ([]string, error) {
	rows, err := a.db.Query(`SELECT title FROM books WHERE authors_id = $1`, id)
	if err != nil {
		a.logger.Error("select for get for authors books failed", "id", id, "error", err.Error())
		return nil, err
	}

	var books []string
	for rows.Next() {
		var book model.Book
		err = rows.Scan(&book.Title)
		if err != nil {
			a.logger.Error("scan rows failed for authors books", "id", id, "error", err.Error())
			return nil, err
		}

		books = append(books, book.Title)
	}

	return books, nil
}
