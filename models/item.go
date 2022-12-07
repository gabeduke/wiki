package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Item is used by pop to map your items database table to your go code.
type Item struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	Name          string        `json:"name" db:"name"`
	Description   nulls.String  `json:"description" db:"description"`
	Claimed       bool          `json:"claimed" db:"claimed"`
	Tags          []Tag         `json:"tags" db:"tags" many_to_many:"items_tags"`
	UserWorkspace UserWorkspace `json:"user_workspace" db:"user_workspace" belongs_to:"user_workspace"`
}

// String is not required by pop and may be deleted
func (i Item) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Items is not required by pop and may be deleted
type Items []Item

// String is not required by pop and may be deleted
func (i Items) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Item) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Item) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Item) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
