package users

import (
	"testing"
)

var (
	TestUserCore *Users
	pass         = "HoshimachiSuisei"
	email        = "hoshimachisuisei@hololive.com"
	id           = 1
	username     = "Suisei"
)

func init() {
	TestUserCore = New("SuiseiKawaii")
}

func TestUsers_HashPassword(t *testing.T) {
	type fields struct {
		secretKey string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		Hashed  bool
		wantErr bool
	}{
		{
			name: "hash with secret key",
			fields: fields{
				secretKey: "ZetaPon",
			},
			args: args{
				password: pass,
			},
			Hashed:  true,
			wantErr: false,
		},
		{
			name: "hash with no secret key",
			fields: fields{
				secretKey: "",
			},
			args: args{
				password: pass,
			},
			Hashed:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Users{
				secretKey: tt.fields.secretKey,
			}
			got, err := u.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != tt.args.password) != tt.Hashed {
				t.Errorf("Users.HashPassword() = %v, want %v", got, tt.Hashed)
			}
		})
	}
}

func TestUsers_CheckPasswordHash(t *testing.T) {
	hashPass, _ := TestUserCore.HashPassword(pass)
	falsePass, _ := TestUserCore.HashPassword("WrongPassword")
	type fields struct {
		secretKey string
	}
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "correct hash password",
			fields: fields{
				secretKey: TestUserCore.secretKey,
			},
			args: args{
				pass,
				hashPass,
			},
			want: true,
		},

		{
			name: "wrong password",
			fields: fields{
				secretKey: TestUserCore.secretKey,
			},
			args: args{
				"SoraCute",
				falsePass,
			},
			want: false,
		},

		{
			name: "wrong hass password",
			fields: fields{
				secretKey: TestUserCore.secretKey,
			},
			args: args{
				pass,
				falsePass,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Users{
				secretKey: tt.fields.secretKey,
			}
			if got := u.CheckPasswordHash(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("Users.CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
