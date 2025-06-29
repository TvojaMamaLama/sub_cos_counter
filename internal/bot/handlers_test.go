package bot

import (
	"sub-cos-counter/internal/models"
	"testing"
)

// Mock context for testing
type mockContext struct {
	callbackData string
	senderID     int64
	textContent  string
}

func (m *mockContext) Data() string {
	return m.callbackData
}

func (m *mockContext) Sender() *mockUser {
	return &mockUser{id: m.senderID}
}

func (m *mockContext) Text() string {
	return m.textContent
}

func (m *mockContext) Edit(text string, markup interface{}, opts ...interface{}) error {
	return nil
}

func (m *mockContext) Send(text string, opts ...interface{}) error {
	return nil
}

type mockUser struct {
	id int64
}

func (m *mockUser) ID() int64 {
	return m.id
}

// Test bot state management
func TestBotStateManagement(t *testing.T) {
	// Create a bot instance for testing
	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	userID := int64(12345)

	// Test initial state
	state := bot.getUserState(userID)
	if state.State != StateIdle {
		t.Errorf("Expected initial state to be %s, got %s", StateIdle, state.State)
	}

	// Test setting state
	bot.setState(userID, StateAddingSubscription)
	state = bot.getUserState(userID)
	if state.State != StateAddingSubscription {
		t.Errorf("Expected state to be %s, got %s", StateAddingSubscription, state.State)
	}

	// Test setting data
	bot.setData(userID, "category", models.CategoryWork)
	bot.setData(userID, "currency", models.CurrencyUSD)
	bot.setData(userID, "period_days", 30)

	// Test getting data
	category := bot.getData(userID, "category")
	if category != models.CategoryWork {
		t.Errorf("Expected category to be %s, got %v", models.CategoryWork, category)
	}

	currency := bot.getData(userID, "currency")
	if currency != models.CurrencyUSD {
		t.Errorf("Expected currency to be %s, got %v", models.CurrencyUSD, currency)
	}

	periodDays := bot.getData(userID, "period_days")
	if periodDays != 30 {
		t.Errorf("Expected period_days to be 30, got %v", periodDays)
	}

	// Test clear state
	bot.clearUserState(userID)
	state = bot.getUserState(userID)
	if state.State != StateIdle {
		t.Errorf("Expected state after clear to be %s, got %s", StateIdle, state.State)
	}

	// Data should be cleared
	category = bot.getData(userID, "category")
	if category != nil {
		t.Errorf("Expected category to be nil after clear, got %v", category)
	}
}

// Test category parsing
func TestCategoryParsing(t *testing.T) {
	tests := []struct {
		callbackData     string
		expectedCategory models.Category
	}{
		{"cat_work", models.CategoryWork},
		{"cat_entertainment", models.CategoryEntertainment},
		{"cat_education", models.CategoryEducation},
		{"cat_home", models.CategoryHome},
		{"cat_other", models.CategoryOther},
	}

	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	for _, test := range tests {
		t.Run(test.callbackData, func(t *testing.T) {
			userID := int64(12345)

			// Simulate category selection
			var category models.Category
			switch test.callbackData {
			case "cat_entertainment":
				category = models.CategoryEntertainment
			case "cat_work":
				category = models.CategoryWork
			case "cat_education":
				category = models.CategoryEducation
			case "cat_home":
				category = models.CategoryHome
			case "cat_other":
				category = models.CategoryOther
			}

			bot.setData(userID, "category", category)

			// Verify data was set correctly
			storedCategory := bot.getData(userID, "category")
			if storedCategory != test.expectedCategory {
				t.Errorf("For callback %s, expected category %s, got %v",
					test.callbackData, test.expectedCategory, storedCategory)
			}
		})
	}
}

// Test currency parsing
func TestCurrencyParsing(t *testing.T) {
	tests := []struct {
		callbackData     string
		expectedCurrency models.Currency
	}{
		{"curr_usd", models.CurrencyUSD},
		{"curr_rub", models.CurrencyRUB},
	}

	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	for _, test := range tests {
		t.Run(test.callbackData, func(t *testing.T) {
			userID := int64(12345)

			// Simulate currency selection
			var currency models.Currency
			switch test.callbackData {
			case "curr_usd":
				currency = models.CurrencyUSD
			case "curr_rub":
				currency = models.CurrencyRUB
			}

			bot.setData(userID, "currency", currency)

			// Verify data was set correctly
			storedCurrency := bot.getData(userID, "currency")
			if storedCurrency != test.expectedCurrency {
				t.Errorf("For callback %s, expected currency %s, got %v",
					test.callbackData, test.expectedCurrency, storedCurrency)
			}
		})
	}
}

// Test complete subscription flow data preservation
func TestSubscriptionFlowDataPreservation(t *testing.T) {
	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	userID := int64(12345)

	// Simulate complete flow
	bot.setState(userID, StateAddingSubscription)

	// Step 1: Category
	bot.setData(userID, "category", models.CategoryWork)

	// Step 2: Currency
	bot.setData(userID, "currency", models.CurrencyUSD)

	// Step 3: Period
	bot.setData(userID, "period_days", 30)

	// Step 4: Auto renewal
	bot.setData(userID, "auto_renewal", true)

	// Step 5: Name
	bot.setData(userID, "name", "Test Subscription")

	// Step 6: Cost
	bot.setData(userID, "cost", models.Money(1999)) // 19.99

	// Verify all data is preserved
	category := bot.getData(userID, "category")
	if category != models.CategoryWork {
		t.Errorf("Category lost: expected %s, got %v", models.CategoryWork, category)
	}

	currency := bot.getData(userID, "currency")
	if currency != models.CurrencyUSD {
		t.Errorf("Currency lost: expected %s, got %v", models.CurrencyUSD, currency)
	}

	periodDays := bot.getData(userID, "period_days")
	if periodDays != 30 {
		t.Errorf("Period lost: expected 30, got %v", periodDays)
	}

	autoRenewal := bot.getData(userID, "auto_renewal")
	if autoRenewal != true {
		t.Errorf("Auto renewal lost: expected true, got %v", autoRenewal)
	}

	name := bot.getData(userID, "name")
	if name != "Test Subscription" {
		t.Errorf("Name lost: expected 'Test Subscription', got %v", name)
	}

	cost := bot.getData(userID, "cost")
	if cost != models.Money(1999) {
		t.Errorf("Cost lost: expected %d, got %v", models.Money(1999), cost)
	}
}
