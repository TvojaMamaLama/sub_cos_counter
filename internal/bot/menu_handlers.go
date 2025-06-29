package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sub-cos-counter/internal/models"

	"gopkg.in/telebot.v3"
)

func (b *Bot) handleMySubscriptions(c telebot.Context) error {
	ctx := context.Background()
	subscriptions, err := b.subscriptionService.GetAllActiveSubscriptions(ctx)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка получения подписок: %v", err))
	}

	if len(subscriptions) == 0 {
		return c.Edit("📋 *Мои подписки*\n\nУ вас пока нет активных подписок.", &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{btnAddSubscription},
				{btnBack},
			},
		}, telebot.ModeMarkdown)
	}

	text := "📋 *Мои подписки*\n\n"
	keyboard := [][]telebot.InlineButton{}

	for _, sub := range subscriptions {
		currencySymbol := "$"
		if sub.Currency == models.CurrencyRUB {
			currencySymbol = "₽"
		}

		status := ""
		if sub.IsPaymentDue() {
			status = " ⚠️"
		}

		text += fmt.Sprintf("• %s - %s%s%s\n  📅 Следующий платеж: %s\n\n",
			sub.Name, sub.Cost.String(), currencySymbol, status, sub.NextPayment.Format("02.01.2006"))

		// Create action buttons for each subscription
		payBtn := telebot.InlineButton{
			Unique: fmt.Sprintf("pay_%d", sub.ID),
			Text:   fmt.Sprintf("✅ Оплатить %s", sub.Name),
		}
		deleteBtn := telebot.InlineButton{
			Unique: fmt.Sprintf("delete_%d", sub.ID),
			Text:   fmt.Sprintf("❌ Удалить %s", sub.Name),
		}

		// Register handlers for these specific buttons
		b.bot.Handle(&payBtn, b.handlePaySubscription)
		b.bot.Handle(&deleteBtn, b.handleDeleteSubscription)

		keyboard = append(keyboard, []telebot.InlineButton{payBtn})
		keyboard = append(keyboard, []telebot.InlineButton{deleteBtn})
	}

	keyboard = append(keyboard, []telebot.InlineButton{btnBack})

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: keyboard,
	}, telebot.ModeMarkdown)
}

func (b *Bot) handlePaySubscription(c telebot.Context) error {
	data := c.Data()
	idStr := strings.TrimPrefix(data, "pay_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Send("❌ Некорректный ID подписки")
	}

	ctx := context.Background()
	err = b.subscriptionService.MarkAsPaid(ctx, id)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка при отметке об оплате: %v", err))
	}

	subscription, err := b.subscriptionService.GetSubscriptionByID(ctx, id)
	if err != nil {
		return c.Send("❌ Ошибка получения подписки")
	}

	currencySymbol := "$"
	if subscription.Currency == models.CurrencyRUB {
		currencySymbol = "₽"
	}

	text := fmt.Sprintf("✅ *Платеж отмечен!*\n\n"+
		"📝 Подписка: %s\n"+
		"💰 Сумма: %s%s\n"+
		"📅 Следующий платеж: %s",
		subscription.Name,
		subscription.Cost.String(), currencySymbol,
		subscription.NextPayment.Format("02.01.2006"))

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnMySubscriptions},
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleDeleteSubscription(c telebot.Context) error {
	data := c.Data()
	idStr := strings.TrimPrefix(data, "delete_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Send("❌ Некорректный ID подписки")
	}

	ctx := context.Background()
	subscription, err := b.subscriptionService.GetSubscriptionByID(ctx, id)
	if err != nil {
		return c.Send("❌ Ошибка получения подписки")
	}

	err = b.subscriptionService.DeleteSubscription(ctx, id)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка при удалении: %v", err))
	}

	text := fmt.Sprintf("🗑️ *Подписка удалена*\n\n"+
		"📝 Название: %s\n"+
		"Подписка была успешно удалена из списка.",
		subscription.Name)

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnMySubscriptions},
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleMonthlyExpense(c telebot.Context) error {
	ctx := context.Background()

	// Get current month expenses
	currentExpenses, err := b.analyticsService.GetCurrentMonthExpense(ctx)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка получения расходов: %v", err))
	}

	// Get monthly recurring costs
	recurringCosts, err := b.analyticsService.GetMonthlyRecurringCost(ctx)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка расчета месячных расходов: %v", err))
	}

	text := "💰 *Месячные расходы*\n\n"

	// Current month actual expenses
	text += "📊 *Оплачено в этом месяце:*\n"
	if len(currentExpenses) == 0 {
		text += "Пока нет платежей в этом месяце\n\n"
	} else {
		for _, expense := range currentExpenses {
			currencySymbol := "$"
			if expense.Currency == models.CurrencyRUB {
				currencySymbol = "₽"
			}
			text += fmt.Sprintf("• %s%s (%d платежей)\n", expense.TotalAmount.String(), currencySymbol, expense.Count)
		}
		text += "\n"
	}

	// Monthly recurring costs
	text += "🔄 *Ежемесячные расходы (все подписки):*\n"
	if len(recurringCosts) == 0 {
		text += "Нет активных подписок\n"
	} else {
		for currency, amount := range recurringCosts {
			currencySymbol := "$"
			if currency == models.CurrencyRUB {
				currencySymbol = "₽"
			}
			text += fmt.Sprintf("• %s%s в месяц\n", amount.String(), currencySymbol)
		}
	}

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleAnalytics(c telebot.Context) error {
	ctx := context.Background()
	analytics, err := b.analyticsService.GetCurrentMonthCategoryAnalytics(ctx)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка получения аналитики: %v", err))
	}

	text := "📊 *Аналитика по категориям*\n\n"

	if len(analytics) == 0 {
		text += "Нет данных за текущий месяц"
	} else {
		for category, summaries := range analytics {
			text += fmt.Sprintf("%s\n", getCategoryEmoji(category))
			for _, summary := range summaries {
				currencySymbol := "$"
				if summary.Currency == models.CurrencyRUB {
					currencySymbol = "₽"
				}
				text += fmt.Sprintf("  %s%s (%d платежей)\n", summary.TotalAmount.String(), currencySymbol, summary.Count)
			}
			text += "\n"
		}
	}

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleHistory(c telebot.Context) error {
	ctx := context.Background()
	payments, err := b.analyticsService.GetPaymentHistory(ctx, 10) // Last 10 payments
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка получения истории: %v", err))
	}

	text := "📜 *История платежей*\n\n"

	if len(payments) == 0 {
		text += "История платежей пуста"
	} else {
		for _, payment := range payments {
			currencySymbol := "$"
			if payment.Currency == models.CurrencyRUB {
				currencySymbol = "₽"
			}

			statusIcon := "✅"
			switch payment.Status {
			case models.PaymentStatusPending:
				statusIcon = "⏳"
			case models.PaymentStatusFailed:
				statusIcon = "❌"
			}

			text += fmt.Sprintf("%s %s%s - %s\n",
				statusIcon, payment.Amount.String(), currencySymbol, payment.PaidAt.Format("02.01.2006 15:04"))
		}
	}

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}

func (b *Bot) handleSettings(c telebot.Context) error {
	text := "⚙️ *Настройки*\n\nДоступные функции:\n\n" +
		"• Поддержка валют: USD, RUB\n" +
		"• Автоматические уведомления о платежах\n" +
		"• Аналитика по категориям\n" +
		"• История всех операций\n\n" +
		"Для настройки уведомлений обратитесь к разработчику."

	return c.Edit(text, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnBack},
		},
	}, telebot.ModeMarkdown)
}
