package client

import "context"

type ExternalServiceClient interface {
	SendNotification(ctx context.Context, userID string, message string) error
}
