package utils

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ULID type for GORM
type ULID ulid.ULID

// Implement fmt.Stringer interface for ULID
func (u ULID) String() string {
	return ulid.ULID(u).String()
}

// Implement GORM custom scanner & valuer for ULID
func (u *ULID) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan ULID: %v", value)
	}
	ulid, err := ulid.Parse(str)
	if err != nil {
		return err
	}
	*u = ULID(ulid)
	return nil
}

func (u ULID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Generate new ULID before saving to DB
func (u *ULID) BeforeCreate(tx *gorm.DB) error {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	newULID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	*u = ULID(newULID)
	return nil
}
