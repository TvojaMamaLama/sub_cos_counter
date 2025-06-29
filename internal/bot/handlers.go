package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sub-cos-counter/internal/models"
	"time"

	"gopkg.in/telebot.v3"
)

func (b *Bot) setupHandlers() {
	// Start command
	b.bot.Handle("/start", b.handleStart)
	
	// Main menu callbacks
	b.bot.Handle(&btnAddSubscription, b.handleAddSubscription)
	b.bot.Handle(&btnMySubscriptions, b.handleMySubscriptions)
	b.bot.Handle(&btnMonthlyExpense, b.handleMonthlyExpense)
	b.bot.Handle(&btnAnalytics, b.handleAnalytics)
	b.bot.Handle(&btnHistory, b.handleHistory)
	b.bot.Handle(&btnSettings, b.handleSettings)
	
	// Category selection callbacks
	b.bot.Handle(&btnCategoryEntertainment, b.handleCategorySelection)
	b.bot.Handle(&btnCategoryWork, b.handleCategorySelection)
	b.bot.Handle(&btnCategoryEducation, b.handleCategorySelection)
	b.bot.Handle(&btnCategoryHome, b.handleCategorySelection)
	b.bot.Handle(&btnCategoryOther, b.handleCategorySelection)
	
	// Currency selection callbacks
	b.bot.Handle(&btnCurrencyUSD, b.handleCurrencySelection)
	b.bot.Handle(&btnCurrencyRUB, b.handleCurrencySelection)
	
	// Period selection callbacks
	b.bot.Handle(&btnPeriodWeek, b.handlePeriodSelection)
	b.bot.Handle(&btnPeriodMonth, b.handlePeriodSelection)
	b.bot.Handle(&btnPeriodYear, b.handlePeriodSelection)
	b.bot.Handle(&btnPeriodCustom, b.handlePeriodSelection)
	
	// Auto renewal callbacks
	b.bot.Handle(&btnAutoRenewalYes, b.handleAutoRenewalSelection)
	b.bot.Handle(&btnAutoRenewalNo, b.handleAutoRenewalSelection)
	
	// Action callbacks
	b.bot.Handle(&btnBack, b.handleBack)
	
	// Text message handler for states
	b.bot.Handle(telebot.OnText, b.handleTextMessage)
	
	// Callback query handler for dynamic callbacks
	b.bot.Handle(telebot.OnCallback, b.handleCallbackQuery)
}

func (b *Bot) handleStart(c telebot.Context) error {
	return b.showMainMenu(c)
}

func (b *Bot) showMainMenu(c telebot.Context) error {
	text := "🏠 *Главное меню трекера подписок*\n\nВыберите действие:"
	
	return c.Send(text, &telebot.ReplyMarkup{
		InlineKeyboard: mainMenuKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleAddSubscription(c telebot.Context) error {
	userID := c.Sender().ID
	b.setState(userID, StateAddingSubscription)
	b.clearUserState(userID)
	
	text := "📝 *Добавление новой подписки*\n\nВыберите категорию:"
	
	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: categoryKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleCategorySelection(c telebot.Context) error {
	userID := c.Sender().ID
	
	var category models.Category
	switch c.Data() {
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
	
	b.setData(userID, "category", category)
	
	text := fmt.Sprintf("📝 *Добавление подписки*\n\n✅ Категория: %s\n\nВыберите валюту:", getCategoryEmoji(category))
	
	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: currencyKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleCurrencySelection(c telebot.Context) error {
	userID := c.Sender().ID
	
	var currency models.Currency
	currencyText := ""
	switch c.Data() {
	case "curr_usd":
		currency = models.CurrencyUSD
		currencyText = "💵 USD"
	case "curr_rub":
		currency = models.CurrencyRUB
		currencyText = "🔹 RUB"
	}
	
	b.setData(userID, "currency", currency)
	
	category := b.getData(userID, "category").(models.Category)
	
	text := fmt.Sprintf("📝 *Добавление подписки*\n\n✅ Категория: %s\n✅ Валюта: %s\n\nВыберите период оплаты:", 
		getCategoryEmoji(category), currencyText)
	
	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: periodKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handlePeriodSelection(c telebot.Context) error {
	userID := c.Sender().ID
	
	var periodDays int
	periodText := ""
	switch c.Data() {
	case "period_week":
		periodDays = 7
		periodText = "🗓️ Неделя"
	case "period_month":
		periodDays = 30
		periodText = "📅 Месяц"
	case "period_year":
		periodDays = 365
		periodText = "📆 Год"
	case "period_custom":
		b.setState(userID, StateWaitingForDate)
		return c.Edit("📝 Введите количество дней между платежами (например: 14):", &telebot.ReplyMarkup{
			InlineKeyboard: backKeyboard,
		})
	}
	
	b.setData(userID, "period_days", periodDays)
	
	category := b.getData(userID, "category").(models.Category)
	currency := b.getData(userID, "currency").(models.Currency)
	currencyText := ""
	if currency == models.CurrencyUSD {
		currencyText = "💵 USD"
	} else {
		currencyText = "🔹 RUB"
	}
	
	text := fmt.Sprintf("📝 *Добавление подписки*\n\n✅ Категория: %s\n✅ Валюта: %s\n✅ Период: %s\n\nВключить автопродление?", 
		getCategoryEmoji(category), currencyText, periodText)
	
	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: autoRenewalKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleAutoRenewalSelection(c telebot.Context) error {
	userID := c.Sender().ID
	
	autoRenewal := c.Data() == "auto_yes"
	autoRenewalText := "❌ Нет"
	if autoRenewal {
		autoRenewalText = "✅ Да"
	}
	
	b.setData(userID, "auto_renewal", autoRenewal)
	b.setState(userID, StateWaitingForName)
	
	category := b.getData(userID, "category").(models.Category)
	currency := b.getData(userID, "currency").(models.Currency)
	periodDays := b.getData(userID, "period_days").(int)
	
	currencyText := ""
	periodText := ""
	if currency == models.CurrencyUSD {
		currencyText = "💵 USD"
	} else {
		currencyText = "🔹 RUB"
	}
	
	switch periodDays {
	case 7:
		periodText = "🗓️ Неделя"
	case 30:
		periodText = "📅 Месяц"
	case 365:
		periodText = "📆 Год"
	default:
		periodText = fmt.Sprintf("⚡ %d дней", periodDays)
	}
	
	text := fmt.Sprintf("📝 *Добавление подписки*\n\n✅ Категория: %s\n✅ Валюта: %s\n✅ Период: %s\n✅ Автопродление: %s\n\n💬 Введите название подписки:", 
		getCategoryEmoji(category), currencyText, periodText, autoRenewalText)
	
	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: backKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleTextMessage(c telebot.Context) error {
	userID := c.Sender().ID
	state := b.getUserState(userID)
	
	switch state.State {
	case StateWaitingForName:
		return b.handleNameInput(c)
	case StateWaitingForCost:
		return b.handleCostInput(c)
	case StateWaitingForDate:
		return b.handleDateInput(c)
	default:
		return b.showMainMenu(c)
	}
}

func (b *Bot) handleNameInput(c telebot.Context) error {
	userID := c.Sender().ID
	name := strings.TrimSpace(c.Text())
	
	if name == "" {
		return c.Send("❌ Название не может быть пустым. Попробуйте еще раз:")
	}
	
	b.setData(userID, "name", name)
	b.setState(userID, StateWaitingForCost)
	
	currency := b.getData(userID, "currency").(models.Currency)
	currencySymbol := "$"
	if currency == models.CurrencyRUB {
		currencySymbol = "₽"
	}
	
	return c.Send(fmt.Sprintf("💰 Введите стоимость подписки в %s (например: 15.99):", currencySymbol))
}

func (b *Bot) handleCostInput(c telebot.Context) error {
	userID := c.Sender().ID
	costStr := strings.TrimSpace(c.Text())
	
	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil || cost <= 0 {
		return c.Send("❌ Некорректная стоимость. Введите число больше 0 (например: 15.99):")
	}
	
	b.setData(userID, "cost", cost)
	
	// Create subscription
	return b.createSubscription(c)
}

func (b *Bot) handleDateInput(c telebot.Context) error {
	userID := c.Sender().ID
	daysStr := strings.TrimSpace(c.Text())
	
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		return c.Send("❌ Некорректное количество дней. Введите число больше 0:")
	}
	
	b.setData(userID, "period_days", days)
	
	category := b.getData(userID, "category").(models.Category)
	currency := b.getData(userID, "currency").(models.Currency)
	currencyText := ""
	if currency == models.CurrencyUSD {
		currencyText = "💵 USD"
	} else {
		currencyText = "🔹 RUB"
	}
	
	text := fmt.Sprintf("📝 *Добавление подписки*\n\n✅ Категория: %s\n✅ Валюта: %s\n✅ Период: ⚡ %d дней\n\nВключить автопродление?", 
		getCategoryEmoji(category), currencyText, days)
	
	return c.Send(text, &telebot.ReplyMarkup{
		InlineKeyboard: autoRenewalKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) createSubscription(c telebot.Context) error {
	userID := c.Sender().ID
	
	name := b.getData(userID, "name").(string)
	cost := b.getData(userID, "cost").(float64)
	currency := b.getData(userID, "currency").(models.Currency)
	category := b.getData(userID, "category").(models.Category)
	periodDays := b.getData(userID, "period_days").(int)
	autoRenewal := b.getData(userID, "auto_renewal").(bool)
	
	// Set next payment date (starting from today + period)
	nextPayment := time.Now().AddDate(0, 0, periodDays)
	
	req := &models.CreateSubscriptionRequest{
		Name:        name,
		Cost:        cost,
		Currency:    currency,
		PeriodDays:  periodDays,
		NextPayment: nextPayment,
		Category:    category,
		AutoRenewal: autoRenewal,
	}
	
	ctx := context.Background()
	subscription, err := b.subscriptionService.CreateSubscription(ctx, req)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка при создании подписки: %v", err))
	}
	
	// Clear user state
	b.clearUserState(userID)
	
	currencySymbol := "$"
	if currency == models.CurrencyRUB {
		currencySymbol = "₽"
	}
	
	periodText := ""
	switch periodDays {
	case 7:
		periodText = "неделю"
	case 30:
		periodText = "месяц"
	case 365:
		periodText = "год"
	default:
		periodText = fmt.Sprintf("%d дней", periodDays)
	}
	
	text := fmt.Sprintf("✅ *Подписка успешно добавлена!*\n\n"+
		"📝 Название: %s\n"+
		"💰 Стоимость: %.2f%s\n"+
		"📅 Период: каждые %s\n"+
		"🗓️ Следующий платеж: %s\n"+
		"📂 Категория: %s\n"+
		"🔄 Автопродление: %s",
		subscription.Name,
		subscription.Cost, currencySymbol,
		periodText,
		subscription.NextPayment.Format("02.01.2006"),
		getCategoryEmoji(category),
		getBoolEmoji(autoRenewal))
	
	return c.Send(text, &telebot.ReplyMarkup{
		InlineKeyboard: mainMenuKeyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleBack(c telebot.Context) error {
	userID := c.Sender().ID
	b.clearUserState(userID)
	return b.showMainMenu(c)
}

func (b *Bot) handleCallbackQuery(c telebot.Context) error {
	data := c.Data()
	
	// Handle subscription actions
	if strings.HasPrefix(data, "pay_") {
		return b.handlePaySubscription(c)
	}
	if strings.HasPrefix(data, "delete_") {
		return b.handleDeleteSubscription(c)
	}
	
	return nil
}

func getCategoryEmoji(category models.Category) string {
	switch category {
	case models.CategoryEntertainment:
		return "🎮 Развлечения"
	case models.CategoryWork:
		return "💼 Работа"
	case models.CategoryEducation:
		return "📚 Обучение"
	case models.CategoryHome:
		return "🏠 Дом"
	case models.CategoryOther:
		return "📦 Другое"
	default:
		return "📦 Другое"
	}
}

func getBoolEmoji(value bool) string {
	if value {
		return "✅ Да"
	}
	return "❌ Нет"
}