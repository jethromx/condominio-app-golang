package payments

import (
	"time"

	"com.mx/crud/internal/models"
)

type PaymentInput struct {
	ID            uint      `json:"id"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"date"`
	ResidentID    uint      `json:"resident_id"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
}

type PaymentOutput struct {
	ID            uint      `json:"id"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"date"`
	ResidentID    uint      `json:"resident_id"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
}

func MapPaymentInputToModel(input *PaymentInput) *models.Payment {
	return &models.Payment{
		Amount:        input.Amount,
		PaymentDate:   input.PaymentDate,
		ResidentID:    input.ResidentID,
		PaymentMethod: input.PaymentMethod,
		Status:        input.Status,
		Description:   input.Description,
	}
}

func MapPaymentToOutput(payment *models.Payment) *PaymentOutput {
	return &PaymentOutput{
		ID:            payment.ID,
		Amount:        payment.Amount,
		PaymentDate:   payment.PaymentDate,
		ResidentID:    payment.ResidentID,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		Description:   payment.Description,
	}
}

func MapPaymentToOutputList(payments []models.Payment) []*PaymentOutput {
	var output []*PaymentOutput
	for _, payment := range payments {
		output = append(output, MapPaymentToOutput(&payment))
	}
	return output
}
