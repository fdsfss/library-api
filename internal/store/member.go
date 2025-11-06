package store

import (
	"errors"
	"library-api/internal/model"
)

func (m *MemberStore) Get() ([]model.Member, error) {
	rows, err := m.db.Query(`SELECT * FROM members`)
	if err != nil {
		m.logger.Error("get all failed for members", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var member model.Member
		err = rows.Scan(&member.ID, &member.FullName)
		if err != nil {
			m.logger.Error("get all failed for members", "error", err.Error())
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (m *MemberStore) Create(member *model.Member) error {
	_, err := m.db.Exec(`INSERT INTO members(id, full_name) VALUES ($1, $2)`, &member.ID, &member.FullName)
	if err != nil {
		m.logger.Error("create failed for members", "error", err.Error())
		return err
	}

	return nil
}

func (m *MemberStore) Exists(id string) error {
	rows, err := m.db.Query(`SELECT EXISTS (SELECT 1 FROM members WHERE id = $1)`, id)
	if err != nil {
		m.logger.Info("id doesn't exist in members", "info", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			m.logger.Info("scan error", "id", id, "error", err.Error())
			return err
		}

		if !exists {
			m.logger.Info("member does not exist", "id", id)
			return errors.New("member does not exist")
		}
	}

	return nil
}

func (m *MemberStore) Update(id string, member *model.Member) error {
	_, err := m.db.Exec(`UPDATE members SET full_name = $1 WHERE id = $2`, &member.FullName, id)
	if err != nil {
		m.logger.Error("update failed for members", "id", id, "error", err.Error())
		return err
	}

	return nil
}

func (m *MemberStore) Delete(id string) error {
	_, err := m.db.Exec(`DELETE FROM members WHERE id = $1`, id)
	if err != nil {
		m.logger.Error("delete failed for members", "id", id, "error", err.Error())
		return err
	}

	return nil
}
