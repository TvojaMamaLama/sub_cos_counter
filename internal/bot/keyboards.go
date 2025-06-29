package bot

import "gopkg.in/telebot.v3"

// Main menu buttons
var (
	btnAddSubscription  = telebot.InlineButton{Unique: "add_sub", Text: "📝 Добавить подписку"}
	btnMySubscriptions  = telebot.InlineButton{Unique: "my_subs", Text: "📋 Мои подписки"}
	btnMonthlyExpense   = telebot.InlineButton{Unique: "monthly", Text: "💰 Месячные расходы"}
	btnAnalytics        = telebot.InlineButton{Unique: "analytics", Text: "📊 Аналитика"}
	btnHistory          = telebot.InlineButton{Unique: "history", Text: "📜 История платежей"}
	btnSettings         = telebot.InlineButton{Unique: "settings", Text: "⚙️ Настройки"}
)

// Category buttons
var (
	btnCategoryEntertainment = telebot.InlineButton{Unique: "cat_entertainment", Text: "🎮 Развлечения"}
	btnCategoryWork          = telebot.InlineButton{Unique: "cat_work", Text: "💼 Работа"}
	btnCategoryEducation     = telebot.InlineButton{Unique: "cat_education", Text: "📚 Обучение"}
	btnCategoryHome          = telebot.InlineButton{Unique: "cat_home", Text: "🏠 Дом"}
	btnCategoryOther         = telebot.InlineButton{Unique: "cat_other", Text: "📦 Другое"}
)

// Currency buttons
var (
	btnCurrencyUSD = telebot.InlineButton{Unique: "curr_usd", Text: "💵 USD"}
	btnCurrencyRUB = telebot.InlineButton{Unique: "curr_rub", Text: "🔹 RUB"}
)

// Period buttons
var (
	btnPeriodWeek   = telebot.InlineButton{Unique: "period_week", Text: "🗓️ Неделя"}
	btnPeriodMonth  = telebot.InlineButton{Unique: "period_month", Text: "📅 Месяц"}
	btnPeriodYear   = telebot.InlineButton{Unique: "period_year", Text: "📆 Год"}
	btnPeriodCustom = telebot.InlineButton{Unique: "period_custom", Text: "⚡ Другое"}
)

// Auto renewal buttons
var (
	btnAutoRenewalYes = telebot.InlineButton{Unique: "auto_yes", Text: "✅ Да"}
	btnAutoRenewalNo  = telebot.InlineButton{Unique: "auto_no", Text: "❌ Нет"}
)

// Navigation buttons
var (
	btnBack = telebot.InlineButton{Unique: "back", Text: "⬅️ Назад"}
)

// Keyboards
var mainMenuKeyboard = [][]telebot.InlineButton{
	{btnAddSubscription, btnMySubscriptions},
	{btnMonthlyExpense, btnAnalytics},
	{btnHistory, btnSettings},
}

var categoryKeyboard = [][]telebot.InlineButton{
	{btnCategoryEntertainment, btnCategoryWork},
	{btnCategoryEducation, btnCategoryHome},
	{btnCategoryOther},
	{btnBack},
}

var currencyKeyboard = [][]telebot.InlineButton{
	{btnCurrencyUSD, btnCurrencyRUB},
	{btnBack},
}

var periodKeyboard = [][]telebot.InlineButton{
	{btnPeriodWeek, btnPeriodMonth},
	{btnPeriodYear, btnPeriodCustom},
	{btnBack},
}

var autoRenewalKeyboard = [][]telebot.InlineButton{
	{btnAutoRenewalYes, btnAutoRenewalNo},
	{btnBack},
}

var backKeyboard = [][]telebot.InlineButton{
	{btnBack},
}