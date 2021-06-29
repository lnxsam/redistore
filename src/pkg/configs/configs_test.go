package configs_test

import (
	"log"
	"redistore/pkg/configs"
	"testing"

	"github.com/joho/godotenv"
)

func TestEnv(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get_from_env",
			args: args{
				key: "APP_NAME",
			},
			want: "Redistore",
		},
	}

	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal("Error loading .env file in configs test", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := configs.Env(tt.args.key); got != tt.want {
				t.Errorf("Env() = %v, want %v", got, tt.want)
			}
		})
	}
}
