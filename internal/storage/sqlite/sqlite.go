package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"cards/internal/domain/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAll(ctx context.Context, userUID uint64) ([]models.Card, error) {
	const op = "storage.sqlite.GetAll"

	stmt, err := s.db.Prepare("SELECT id, number, cvv, month, year, info FROM cards WHERE user_uid = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cards []models.Card

	for rows.Next() {
		var card models.Card
		err = rows.Scan(&card.ID, &card.Number, &card.CVV, &card.Month, &card.Year, &card.Info)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (s *Storage) SaveCard(
	ctx context.Context,
	number string,
	cvv string,
	month string,
	year string,
	info string,
	userUID uint64,
) (uint64, error) {
	const op = "storage.sqlite.SaveCard"

	stmt, err := s.db.Prepare("INSERT INTO cards(number, cvv, month, year, info, user_uid) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, number, cvv, month, year, info, userUID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uint64(id), nil
}

func (s *Storage) UpdateCard(
	ctx context.Context,
	id uint64,
	number string,
	cvv string,
	month string,
	year string,
	info string,
	userUID uint64,
) error {
	const op = "storage.sqlite.UpdateCard"

	stmt, err := s.db.Prepare("UPDATE cards SET number = ?, cvv = ?, month = ?, year = ?, info = ? WHERE id = ? AND user_uid = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, number, cvv, month, year, info, id, userUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
