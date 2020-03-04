package untraceable

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Leontel represents the entity model to retrieve from db
type Leontel struct {
	LeaID int64   `sql:"column:lea_id"`
	Phone *string `sql:"column:TELEFONO"`
	SouID int64   `sql:"column:lea_source"`
}

// Untraceable represents the entity model of untraceable leads
type Untraceable struct {
	gorm.Model
	LeaID   int64 `sql:"column:lea_id"`
	Phone   *string
	SouID   int64
	DDI     string
	SmsDate time.Time `sql:"DEFAULT:current_timestamp, sms_date"`
}

// Candidates is an struct to handle sms info for each campaign
type Candidates struct {
	Desc  string
	DDI   string
	Leads []Untraceable
}

// Handler is a struct to handle the data used in this functionality
type Handler struct {
	Storer     Storer
	LLeidanet  LLeidanet
	Candidates []Candidates
	Leads      []Untraceable
	Errors     []error
}

// TableName sets the default table name
func (Untraceable) TableName() string {
	return "untraceable"
}
