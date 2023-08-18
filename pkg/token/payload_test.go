package token

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestNewPayload(t *testing.T) {
	type args struct {
		Username string
		UserID   string
		Duration time.Duration
	}
	arg := args{
		Username: "Username",
		UserID:   "1",
		Duration: 15 * time.Minute,
	}

	id, err := strconv.ParseInt(arg.UserID, 10, 64)
	if err != nil {
		log.Fatal("Error when parsing user id")
	}
	payload, err := NewPayload(arg.Username, id, arg.Duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, id, payload.UserID)
	require.Equal(t, arg.Username, payload.Username)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, 100*time.Millisecond)
	require.WithinDuration(t, time.Now().Add(arg.Duration), payload.ExpiredAt, 100*time.Millisecond)
}

func TestPayload_Valid(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		UserID    int64
		Username  string
		IssuedAt  time.Time
		ExpiredAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Does not expired", fields{ExpiredAt: time.Now().Add(15 * time.Minute)}, false},
		{"Already expired", fields{ExpiredAt: time.Now().Add(-15 * time.Minute)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &Payload{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Username:  tt.fields.Username,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiredAt: tt.fields.ExpiredAt,
			}
			if err := payload.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
