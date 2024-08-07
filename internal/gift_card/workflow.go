package gift_card

import (
	"time"

	"github.com/milovidov983/oms-temporal/pkg/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func IssueGiftCard(ctx workflow.Context, input models.GiftCardOrderRequest) (string, error) {
	// Настраиваем политику ретраев
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        42, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}
	// настраиваем опции запуска всех активити
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Оплачиваем карту
	var paymentTransactionID string
	paymentErr := workflow.ExecuteActivity(ctx, Pay, input).Get(ctx, &paymentTransactionID)

	if paymentErr != nil {

		// Все плохо не получилось оплатить. Шлём пуш что не дадим карту
		notificationErr := workflow.ExecuteActivity(ctx, SendFailureNotification, input).Get(ctx, nil)
		if notificationErr != nil {
			return "", notificationErr
		}

		// Вероятно надо дёрнуть особый callback на такой случай
		// чтобы сайтик смог отобразить информацию о проблеме
		// call callback etc...

		return "", paymentErr
	}

	// Создаём подарочную карту
	var giftCardNumber string
	creatingGiftCardErr := workflow.ExecuteActivity(ctx, GetGiftCardNumber, input).Get(ctx, &giftCardNumber)

	if creatingGiftCardErr != nil {

		// Все плохо не получилось создать подарочную карту
		// Пытаемся вернуть бабки
		refundErr := workflow.ExecuteActivity(ctx, Refund, paymentTransactionID, input.Metadata.IdempotencyToken).Get(ctx, nil)
		if refundErr != nil {
			// Все плохо не получилось вернуть бабки
			// Жалуемся в саппорт о том, что у нас лапки
			workflow.ExecuteActivity(ctx, SendSupportAlert, input, "Failed to refund transaction! Please solve this problem!").Get(ctx, nil)
			return "", refundErr
		}

		// Кажется что вернём бабки
		// Клянёмся что скоро вернем бабки
		notificationErr := workflow.ExecuteActivity(ctx, SendRefundNotification, input).Get(ctx, nil)
		if notificationErr != nil {
			// Не получилось уведомить пользователя об ошибке
			return "", notificationErr
		}

		return "", creatingGiftCardErr
	}

	// Отправляем нотификацию что подарочная карта готова
	notificationErr := workflow.ExecuteActivity(ctx, SendSuccessNotification, input).Get(ctx, nil)
	if notificationErr != nil {
		// Не получилось уведомить пользователя о том что карта готова
		// но это не крит, пытаемся вызвать callback и отдать номер подарочной карты

		// log error
	}

	// Возвращаем номер подарочной карты через callback сайтика
	callbackErr := workflow.ExecuteActivity(ctx, ExecuteWebsiteCallback, input, giftCardNumber).Get(ctx, nil)

	if callbackErr != nil {
		// Не получилось вызвать callback. Всё плохо, чел не получит свою карту :(
		// тут можно либо вернуть деньги, но вдруг номер карты всё таки как-то передан
		// поэтому лучше дёрнем наш любимый саппорт
		workflow.ExecuteActivity(ctx, SendSupportAlert, input, "Failed to return a gift card number using a callback!").Get(ctx, nil)

		return "", callbackErr
	}

	return "", nil
}
