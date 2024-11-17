package repository

import (
	"casino-back/internal/app/logger"
	"casino-back/internal/app/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SlotRepository struct {
	db *sqlx.DB
}

func NewSlotRepository(db *sqlx.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) GetSlot(slotName string) (int, error) {
	var id int

	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", "slots")
	err := r.db.Get(&id, query, slotName)

	if err != nil {
		logger.InfoKV("repository error", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *SlotRepository) GetSlotData(slotName string) (model.Slot, error) {
	var slot model.Slot

	query := fmt.Sprintf("SELECT * FROM %s WHERE name = $1", "slots")
	err := r.db.Get(&slot, query, slotName)

	if err != nil {
		logger.InfoKV("repository error", "err", err)
		return model.Slot{}, err
	}

	return slot, nil
}
