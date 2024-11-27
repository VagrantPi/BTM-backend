package model

import (
	"time"

	"github.com/google/uuid"
)

type VerificationType string

const (
	VerificationTypeAutomatic VerificationType = "automatic"
)

type Customer struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key"`
	Phone   string    `gorm:"type:text"`
	PhoneAt time.Time `gorm:"type:timestamptz"`

	// ID Card Related
	IDCardDataNumber     *string    `gorm:"column:id_card_data_number;type:text"`
	IDCardDataExpiration *time.Time `gorm:"column:id_card_data_expiration;type:date"`
	IDCardData           *string    `gorm:"column:id_card_data;type:json"`
	IDCardDataAt         *time.Time `gorm:"column:id_card_data_at;type:timestamptz"`
	IDCardDataRaw        *string    `gorm:"column:id_card_data_raw;type:text"`

	// Basic Info
	Name    *string    `gorm:"type:text"`
	Address *string    `gorm:"type:text"`
	Email   *string    `gorm:"type:text"`
	EmailAt *time.Time `gorm:"type:timestamptz"`

	// Verification Status
	Sanctions    *bool      `gorm:"type:bool"`
	SanctionsAt  *time.Time `gorm:"type:timestamptz"`
	AuthorizedAt *time.Time `gorm:"type:timestamptz"`

	// Photos and Camera
	FrontCameraPath *string    `gorm:"type:text"`
	FrontCameraAt   *time.Time `gorm:"type:timestamptz"`
	IDCardPhotoPath *string    `gorm:"type:text"`
	IDCardPhotoAt   *time.Time `gorm:"type:timestamptz"`

	// US Specific
	USSSN   *string    `gorm:"column:us_ssn;type:text"`
	USSSnAt *time.Time `gorm:"column:us_ssn_at;type:timestamptz"`

	// Override Settings
	SMSOverride           VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	SMSOverrideAt         *time.Time       `gorm:"type:timestamptz"`
	SMSOverrideBy         *uuid.UUID       `gorm:"type:uuid"`
	PhoneOverride         VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	PhoneOverrideAt       *time.Time       `gorm:"type:timestamptz"`
	PhoneOverrideBy       *uuid.UUID       `gorm:"type:uuid"`
	IDCardDataOverride    VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	IDCardDataOverrideAt  *time.Time       `gorm:"type:timestamptz"`
	IDCardDataOverrideBy  *uuid.UUID       `gorm:"type:uuid"`
	IDCardPhotoOverride   VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	IDCardPhotoOverrideAt *time.Time       `gorm:"type:timestamptz"`
	IDCardPhotoOverrideBy *uuid.UUID       `gorm:"type:uuid"`
	FrontCameraOverride   VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	FrontCameraOverrideAt *time.Time       `gorm:"type:timestamptz"`
	FrontCameraOverrideBy *uuid.UUID       `gorm:"type:uuid"`
	SanctionsOverride     VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	SanctionsOverrideAt   *time.Time       `gorm:"type:timestamptz"`
	SanctionsOverrideBy   *uuid.UUID       `gorm:"type:uuid"`
	AuthorizedOverride    VerificationType `gorm:"type:verification_type;not null;default:'automatic'"`
	AuthorizedOverrideAt  *time.Time       `gorm:"type:timestamptz"`
	AuthorizedOverrideBy  *uuid.UUID       `gorm:"type:uuid"`
	USSSnOverride         VerificationType `gorm:"column:us_ssn_override;type:verification_type;not null;default:'automatic'"`
	USSSnOverrideAt       *time.Time       `gorm:"column:us_ssn_override_at;type:timestamptz"`
	USSSnOverrideBy       *uuid.UUID       `gorm:"column:us_ssn_override_by;type:uuid"`

	// Additional Info
	SubscriberInfo   *string    `gorm:"type:json"`
	SubscriberInfoAt *time.Time `gorm:"type:timestamptz"`
	SubscriberInfoBy *uuid.UUID `gorm:"type:uuid"`
	IsTestCustomer   bool       `gorm:"not null;default:false"`

	// Usage Info
	LastAuthAttempt *time.Time `gorm:"type:timestamptz"`
	LastUsedMachine *string    `gorm:"type:text"`
	SuspendedUntil  *time.Time `gorm:"type:timestamptz"`

	// Timestamps
	Created time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// Foreign Key Relationships
	// SMSOverrideByUser         *User `gorm:"foreignKey:SMSOverrideBy"`
	// IDCardDataOverrideByUser  *User `gorm:"foreignKey:IDCardDataOverrideBy"`
	// IDCardPhotoOverrideByUser *User `gorm:"foreignKey:IDCardPhotoOverrideBy"`
	// FrontCameraOverrideByUser *User `gorm:"foreignKey:FrontCameraOverrideBy"`
	// SanctionsOverrideByUser   *User `gorm:"foreignKey:SanctionsOverrideBy"`
	// AuthorizedOverrideByUser  *User `gorm:"foreignKey:AuthorizedOverrideBy"`
	// USSSnOverrideByUser       *User `gorm:"foreignKey:USSSnOverrideBy"`
	// SubscriberInfoByUser      *User `gorm:"foreignKey:SubscriberInfoBy"`
}
