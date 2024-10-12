package email

import (
	"context"
	"log"
	"medods_test_task/service"
)

type MockEmailRepository struct{}

var _ service.EmailRepository = &MockEmailRepository{}

func (em *MockEmailRepository) NotifyUser(ctx context.Context, email string) error {
	log.Printf(" [Warning] Different ip address for %s", email)
	return nil
}
