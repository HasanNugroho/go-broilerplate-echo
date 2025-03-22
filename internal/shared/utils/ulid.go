package utils

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULID type untuk GORM
type ULID struct {
	ulid.ULID
}

// Convert ULID ke String
func (u ULID) String() string {
	return u.ULID.String()
}

// Scan ULID dari database (membaca nilai dari DB)
func (u *ULID) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("ULID value is nil")
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan ULID: %v", value)
	}

	parsedULID, err := ulid.Parse(str)
	if err != nil {
		return err
	}

	u.ULID = parsedULID
	return nil
}

// Menyimpan ULID sebagai string di database
func (u ULID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Generate ULID Baru
func NewULID() ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	newULID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return ULID{newULID}
}
