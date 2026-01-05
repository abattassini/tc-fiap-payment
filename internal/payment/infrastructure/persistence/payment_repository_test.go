package persistence_test

import (
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&entities.Payment{})
	assert.NoError(t, err)

	return db
}

func TestPaymentRepository_AddPayment(t *testing.T) {
	// GIVEN a test database and repository
	db := setupTestDB(t)
	repo := persistence.NewPaymentRepositoryImpl(db)

	payment := &entities.Payment{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}

	// WHEN adding payment
	result, err := repo.AddPayment(payment)

	// THEN payment should be added
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID)
	assert.Equal(t, payment.OrderId, result.OrderId)
	assert.Equal(t, payment.Total, result.Total)
}

func TestPaymentRepository_GetPaymentByOrderId(t *testing.T) {
	// GIVEN a test database with a payment
	db := setupTestDB(t)
	repo := persistence.NewPaymentRepositoryImpl(db)

	payment := &entities.Payment{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}
	repo.AddPayment(payment)

	// WHEN getting payment by order ID
	result, err := repo.GetPaymentByOrderId(1)

	// THEN payment should be returned
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, payment.OrderId, result.OrderId)
	assert.Equal(t, payment.Total, result.Total)
}

func TestPaymentRepository_GetPaymentByOrderId_NotFound(t *testing.T) {
	// GIVEN a test database without payments
	db := setupTestDB(t)
	repo := persistence.NewPaymentRepositoryImpl(db)

	// WHEN getting non-existent payment
	result, err := repo.GetPaymentByOrderId(999)

	// THEN error should be returned
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPaymentRepository_UpdatePayment(t *testing.T) {
	// GIVEN a test database with a payment
	db := setupTestDB(t)
	repo := persistence.NewPaymentRepositoryImpl(db)

	payment := &entities.Payment{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}
	repo.AddPayment(payment)

	// WHEN updating payment
	payment.Status = "Approved"
	err := repo.UpdatePayment(payment)

	// THEN payment should be updated
	assert.NoError(t, err)

	// Verify update
	updated, _ := repo.GetPaymentByOrderId(1)
	assert.Equal(t, "Approved", updated.Status)
}
