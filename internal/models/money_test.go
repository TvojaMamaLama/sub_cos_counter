package models

import (
	"testing"
)

func TestParseMoney(t *testing.T) {
	tests := []struct {
		input    string
		expected Money
		hasError bool
	}{
		// Valid inputs
		{"15.99", Money(1599), false},
		{"15.9", Money(1590), false},
		{"15", Money(1500), false},
		{"0.99", Money(99), false},
		{"0.9", Money(90), false},
		{"0", Money(0), false},
		{"100.00", Money(10000), false},

		// Comma instead of dot
		{"15,99", Money(1599), false},
		{"15,9", Money(1590), false},

		// With spaces
		{" 15.99 ", Money(1599), false},

		// Invalid inputs - too many decimal places
		{"15.999", Money(0), true},
		{"15.9999", Money(0), true},

		// Invalid inputs - non-numeric
		{"abc", Money(0), true},
		{"15.abc", Money(0), true},
		{"", Money(0), true},

		// Invalid inputs - multiple dots
		{"15.99.00", Money(0), true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseMoney(test.input)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", test.input, err)
				}
				if result != test.expected {
					t.Errorf("For input %q, expected %d, got %d", test.input, test.expected, result)
				}
			}
		})
	}
}

func TestMoneyString(t *testing.T) {
	tests := []struct {
		money    Money
		expected string
	}{
		{Money(1599), "15.99"},
		{Money(1590), "15.90"},
		{Money(1500), "15.00"},
		{Money(99), "0.99"},
		{Money(90), "0.90"},
		{Money(0), "0.00"},
		{Money(10000), "100.00"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.money.String()
			if result != test.expected {
				t.Errorf("For money %d, expected %q, got %q", test.money, test.expected, result)
			}
		})
	}
}

func TestMoneyMethods(t *testing.T) {
	m := Money(1599) // 15.99

	// Test Cents()
	if m.Cents() != 1599 {
		t.Errorf("Expected Cents() to return 1599, got %d", m.Cents())
	}

	// Test Dollars()
	if m.Dollars() != 15 {
		t.Errorf("Expected Dollars() to return 15, got %d", m.Dollars())
	}

	// Test CentsOnly()
	if m.CentsOnly() != 99 {
		t.Errorf("Expected CentsOnly() to return 99, got %d", m.CentsOnly())
	}

	// Test IsPositive()
	if !m.IsPositive() {
		t.Error("Expected IsPositive() to return true for 15.99")
	}

	if Money(0).IsPositive() {
		t.Error("Expected IsPositive() to return false for 0")
	}

	// Test IsZero()
	if m.IsZero() {
		t.Error("Expected IsZero() to return false for 15.99")
	}

	if !Money(0).IsZero() {
		t.Error("Expected IsZero() to return true for 0")
	}

	// Test Add()
	result := m.Add(Money(100)) // Add 1.00
	expected := Money(1699)     // 16.99
	if result != expected {
		t.Errorf("Expected Add() to return %d, got %d", expected, result)
	}
}

func TestNewMoneyFromDollarsAndCents(t *testing.T) {
	tests := []struct {
		dollars  int
		cents    int
		expected Money
	}{
		{15, 99, Money(1599)},
		{15, 90, Money(1590)},
		{0, 99, Money(99)},
		{100, 0, Money(10000)},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := NewMoneyFromDollarsAndCents(test.dollars, test.cents)
			if result != test.expected {
				t.Errorf("For %d dollars and %d cents, expected %d, got %d",
					test.dollars, test.cents, test.expected, result)
			}
		})
	}
}
