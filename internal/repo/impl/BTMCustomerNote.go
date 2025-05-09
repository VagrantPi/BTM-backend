package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *repository) CreateCustomerNote(db *gorm.DB, note domain.BTMCustomerNote) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	modelNote := CustomerNotesDomainToModel(note)
	return db.Create(&modelNote).Error
}

func (repo *repository) GetCustomerNotes(db *gorm.DB, customerId uuid.UUID, noteType domain.CustomerNoteType, limit int, page int) ([]domain.BTMCustomerNote, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	sql := db.Model(&model.BTMCustomerNote{}).
		Where("customer_id = ? AND note_type = ?", customerId, noteType.Int())

	var total int64
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 取得資料
	var modelNotes []model.BTMCustomerNote
	if err := sql.
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&modelNotes).Error; err != nil {
		return nil, 0, err
	}

	// 轉換為 domain 對象
	notes := make([]domain.BTMCustomerNote, len(modelNotes))
	for i, modelNote := range modelNotes {
		notes[i] = CustomerNotesModelToDomain(modelNote)
	}

	return notes, total, nil
}

func (repo *repository) GetCustomerNote(db *gorm.DB, noteId uint) (domain.BTMCustomerNote, error) {
	if db == nil {
		return domain.BTMCustomerNote{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var modelNote model.BTMCustomerNote
	if err := db.First(&modelNote, noteId).Error; err != nil {
		return domain.BTMCustomerNote{}, err
	}

	return CustomerNotesModelToDomain(modelNote), nil
}

func (repo *repository) DeleteCustomerNote(db *gorm.DB, noteId uint) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Delete(&model.BTMCustomerNote{}, noteId).Error
}

func (repo *repository) UpdateCustomerNote(db *gorm.DB, note domain.BTMCustomerNote) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Model(&model.BTMCustomerNote{}).
		Where("id = ?", note.ID).
		Updates(CustomerNotesDomainToModel(note)).Error
}

func CustomerNotesDomainToModel(note domain.BTMCustomerNote) model.BTMCustomerNote {
	return model.BTMCustomerNote{
		CustomerId:        note.CustomerId,
		Note:              note.Note,
		OperationUserId:   note.OperationUserId,
		OperationUserName: note.OperationUserName,
		NoteType:          note.NoteType.Int(),
	}
}

func CustomerNotesModelToDomain(note model.BTMCustomerNote) domain.BTMCustomerNote {
	return domain.BTMCustomerNote{
		ID:                uint(note.ID),
		CreatedAt:         note.CreatedAt,
		CustomerId:        note.CustomerId,
		Note:              note.Note,
		OperationUserId:   note.OperationUserId,
		OperationUserName: note.OperationUserName,
		NoteType:          domain.CustomerNoteType(note.NoteType),
	}
}
