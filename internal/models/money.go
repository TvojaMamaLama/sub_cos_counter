package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Money represents a monetary amount stored as cents/kopecks (smallest unit)
type Money int

// NewMoney creates a new Money instance from cents
func NewMoney(cents int) Money {
	return Money(cents)
}

// NewMoneyFromDollarsAndCents creates Money from dollars and cents
func NewMoneyFromDollarsAndCents(dollars, cents int) Money {
	return Money(dollars*100 + cents)
}

// ParseMoney parses a string like "15.99" into Money (stored as cents)
func ParseMoney(str string) (Money, error) {
	str = strings.TrimSpace(str)

	// Replace comma with dot for international support
	str = strings.Replace(str, ",", ".", 1)

	parts := strings.Split(str, ".")
	if len(parts) > 2 {
		return Money(0), fmt.Errorf("invalid money format: %s", str)
	}

	dollars, err := strconv.Atoi(parts[0])
	if err != nil {
		return Money(0), fmt.Errorf("invalid dollars: %s", parts[0])
	}

	cents := 0
	if len(parts) == 2 {
		centStr := parts[1]
		// Validate cents format
		if len(centStr) > 2 {
			return Money(0), fmt.Errorf("too many digits after decimal point: %s (max 2 digits)", centStr)
		}
		// Pad to 2 digits if only 1 digit provided
		if len(centStr) == 1 {
			centStr += "0"
		}

		cents, err = strconv.Atoi(centStr)
		if err != nil {
			return Money(0), fmt.Errorf("invalid cents: %s", centStr)
		}
	}

	return NewMoneyFromDollarsAndCents(dollars, cents), nil
}

// String returns the money as a formatted string like "15.99"
func (m Money) String() string {
	dollars := int(m) / 100
	cents := int(m) % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}

// Cents returns the total amount in cents
func (m Money) Cents() int {
	return int(m)
}

// Dollars returns just the dollar part
func (m Money) Dollars() int {
	return int(m) / 100
}

// CentsOnly returns just the cents part (0-99)
func (m Money) CentsOnly() int {
	return int(m) % 100
}

// Add adds two Money amounts
func (m Money) Add(other Money) Money {
	return Money(int(m) + int(other))
}

// IsZero returns true if the amount is zero
func (m Money) IsZero() bool {
	return m == 0
}

// IsPositive returns true if the amount is positive
func (m Money) IsPositive() bool {
	return m > 0
}

// Value implements the driver.Valuer interface for database storage
func (m Money) Value() (driver.Value, error) {
	// Store as integer (cents)
	return int(m), nil
}

// Scan implements the sql.Scanner interface for database reading
func (m *Money) Scan(value interface{}) error {
	if value == nil {
		*m = Money(0)
		return nil
	}

	switch v := value.(type) {
	case int64:
		*m = Money(v)
	case int32:
		*m = Money(v)
	case int:
		*m = Money(v)
	default:
		return fmt.Errorf("cannot scan %T into Money", value)
	}

	return nil
}
