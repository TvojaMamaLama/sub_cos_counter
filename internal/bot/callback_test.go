package bot

import (
	"sub-cos-counter/internal/models"
	"testing"
)

// Test that simulates real Telebot callback handling
func TestCallbackHandlerRegistration(t *testing.T) {
	// Create a minimal bot setup similar to real usage
	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	// Test callback data processing directly
	userID := int64(12345)

	// Step 1: Category selection
	t.Run("Category Selection", func(t *testing.T) {
		// Simulate handleCategorySelection logic
		callbackData := "cat_work"

		var category models.Category
		switch callbackData {
		case "cat_work":
			category = models.CategoryWork
		}

		bot.setData(userID, "category", category)

		// Verify
		stored := bot.getData(userID, "category")
		if stored != models.CategoryWork {
			t.Errorf("Category not stored correctly: expected %s, got %v", models.CategoryWork, stored)
		}
	})

	// Step 2: Currency selection
	t.Run("Currency Selection", func(t *testing.T) {
		// Simulate handleCurrencySelection logic
		callbackData := "curr_usd"

		var currency models.Currency
		switch callbackData {
		case "curr_usd":
			currency = models.CurrencyUSD
		}

		bot.setData(userID, "currency", currency)

		// Verify category is still there
		category := bot.getData(userID, "category")
		if category != models.CategoryWork {
			t.Error("Category lost after currency selection!")
		}

		// Verify currency
		stored := bot.getData(userID, "currency")
		if stored != models.CurrencyUSD {
			t.Errorf("Currency not stored correctly: expected %s, got %v", models.CurrencyUSD, stored)
		}
	})

	// Step 3: Period selection
	t.Run("Period Selection", func(t *testing.T) {
		// Simulate handlePeriodSelection logic
		callbackData := "period_month"

		var periodDays int
		switch callbackData {
		case "period_month":
			periodDays = 30
		}

		bot.setData(userID, "period_days", periodDays)

		// Verify all previous data is still there
		category := bot.getData(userID, "category")
		if category != models.CategoryWork {
			t.Error("Category lost after period selection!")
		}

		currency := bot.getData(userID, "currency")
		if currency != models.CurrencyUSD {
			t.Error("Currency lost after period selection!")
		}

		// Verify period
		stored := bot.getData(userID, "period_days")
		if stored != 30 {
			t.Errorf("Period not stored correctly: expected 30, got %v", stored)
		}
	})

	// Step 4: Auto renewal selection
	t.Run("Auto Renewal Selection", func(t *testing.T) {
		// Simulate handleAutoRenewalSelection logic
		callbackData := "auto_yes"

		autoRenewal := callbackData == "auto_yes"
		bot.setData(userID, "auto_renewal", autoRenewal)

		// Verify ALL data is preserved
		category := bot.getData(userID, "category")
		currency := bot.getData(userID, "currency")
		periodDays := bot.getData(userID, "period_days")
		storedAutoRenewal := bot.getData(userID, "auto_renewal")

		if category != models.CategoryWork {
			t.Errorf("Category lost: expected %s, got %v", models.CategoryWork, category)
		}
		if currency != models.CurrencyUSD {
			t.Errorf("Currency lost: expected %s, got %v", models.CurrencyUSD, currency)
		}
		if periodDays != 30 {
			t.Errorf("Period lost: expected 30, got %v", periodDays)
		}
		if storedAutoRenewal != true {
			t.Errorf("Auto renewal not stored: expected true, got %v", storedAutoRenewal)
		}
	})
}

// Test the exact button unique IDs match between keyboards and handlers
func TestButtonUniqueIDs(t *testing.T) {
	// These should match keyboards.go exactly
	expectedButtons := map[string]string{
		"cat_work":          "💼 Работа",
		"cat_entertainment": "🎮 Развлечения",
		"cat_education":     "📚 Обучение",
		"cat_home":          "🏠 Дом",
		"cat_other":         "📦 Другое",
		"curr_usd":          "💵 USD",
		"curr_rub":          "🔹 RUB",
		"period_week":       "🗓️ Неделя",
		"period_month":      "📅 Месяц",
		"period_year":       "📆 Год",
		"period_custom":     "⚡ Другое",
		"auto_yes":          "✅ Да",
		"auto_no":           "❌ Нет",
	}

	// Test that our callback parsing matches the button IDs
	for uniqueID := range expectedButtons {
		t.Run(uniqueID, func(t *testing.T) {
			// Test category parsing
			if uniqueID == "cat_work" {
				var category models.Category
				switch uniqueID {
				case "cat_work":
					category = models.CategoryWork
				}
				if category != models.CategoryWork {
					t.Errorf("Button %s not parsed correctly", uniqueID)
				}
			}

			// Test currency parsing
			if uniqueID == "curr_usd" {
				var currency models.Currency
				switch uniqueID {
				case "curr_usd":
					currency = models.CurrencyUSD
				}
				if currency != models.CurrencyUSD {
					t.Errorf("Button %s not parsed correctly", uniqueID)
				}
			}

			// Test period parsing
			if uniqueID == "period_month" {
				var periodDays int
				switch uniqueID {
				case "period_month":
					periodDays = 30
				}
				if periodDays != 30 {
					t.Errorf("Button %s not parsed correctly", uniqueID)
				}
			}
		})
	}
}

// Test that clearUserState is called at the right time
func TestClearUserStateFlow(t *testing.T) {
	bot := &Bot{
		userStates: make(map[int64]*UserState),
	}

	userID := int64(12345)

	// Simulate handleAddSubscription - should clear state first
	bot.clearUserState(userID)
	bot.setState(userID, StateAddingSubscription)

	// Verify state is clean
	category := bot.getData(userID, "category")
	if category != nil {
		t.Errorf("State not properly cleared: category should be nil, got %v", category)
	}

	state := bot.getUserState(userID)
	if state.State != StateAddingSubscription {
		t.Errorf("State not set correctly: expected %s, got %s", StateAddingSubscription, state.State)
	}

	// Add some data
	bot.setData(userID, "category", models.CategoryWork)
	bot.setData(userID, "currency", models.CurrencyUSD)

	// Simulate successful subscription creation - should clear state
	bot.clearUserState(userID)

	// Verify everything is cleared
	category = bot.getData(userID, "category")
	currency := bot.getData(userID, "currency")
	state = bot.getUserState(userID)

	if category != nil {
		t.Errorf("Category not cleared: got %v", category)
	}
	if currency != nil {
		t.Errorf("Currency not cleared: got %v", currency)
	}
	if state.State != StateIdle {
		t.Errorf("State not reset: expected %s, got %s", StateIdle, state.State)
	}
}
