package sign_test

import (
	"hash"
	"testing"

	"github.com/thefrol/kysh-kysh-meow/internal/sign"
)

// проверка подписи для стандартных хешера и енкодера
func TestBytes(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name    string
		hasher  func() hash.Hash
		encoder func([]byte) string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{data: []byte("Выпей ещё чайку"), key: []byte("с этими булочками")},
			want:    "7RWbJONK8sV0rF6W4gJGIa2+Erh9szjEVuzFbDXtddk=",
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got, err := sign.Bytes(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
