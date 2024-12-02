package pgxsql_test

import (
	"context"
	"testing"

	"github.com/mkadit/go-toybox/internal/models"
)

func TestAdapter_CreateUser(t *testing.T) {
	dbTester.CheckParallel(t)
	ad, err := dbTester.GetConnection(context.Background())
	if err != nil {
		t.Error("failed to connect to db: ", err)
	}

	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "insert valid new user",
			args: args{
				user: models.User{
					Email:        "kaela@hololive.com",
					Username:     "kaela",
					PasswordHash: "123123",
				},
			},
			wantErr: false,
		},

		{
			name: "insert duplicate user",
			args: args{
				user: models.User{
					Email:        "kaela@hololive.com",
					Username:     "kaela",
					PasswordHash: "123123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ad.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Adapter.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	_ = ad
}

func TestAdapter_GetUserByEmail(t *testing.T) {
	dbTester.CheckParallel(t)
	ad, err := dbTester.GetConnection(context.Background())
	if err != nil {
		t.Error("failed to connect to db: ", err)
	}

	type args struct {
		email string
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
		wantErr   bool
	}{
		{
			name: "get valid user",
			args: args{
				email: "suisei@hololive.com",
			},
			wantEmpty: false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				gotResult, err := ad.GetUserByEmail(tt.args.email)
				if (err != nil) != tt.wantErr {
					t.Errorf("Adapter.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if (gotResult.ID == 0) != tt.wantEmpty {
					t.Errorf("Adapter.GetUserByEmail() = %v, want %v", gotResult, tt.wantEmpty)
				}

			},
		)
	}
}
