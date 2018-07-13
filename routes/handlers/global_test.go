package handlers

import (
	"testing"
)

func TestParseJWT(t *testing.T) {
	type args struct {
		token  string
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiaXJhbmkifQ.2e1NJw47Gl_DzJcXw5uK1r99Qnm42DRjSYKi2ASFDnQ",
				secret: []byte("your-256-bit-secret"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseJWT(tt.args.token, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
