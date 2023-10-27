package entities

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCrypto(t *testing.T) {
	type args struct {
		shortTitle string
		cost       float64
	}
	tests := []struct {
		name    string
		args    args
		want    *Crypto
		wantErr bool
	}{
		{
			name:    "new crypto create successful",
			args:    args{shortTitle: "ETH", cost: 1.22},
			want:    &Crypto{ShortTitle: "ETH", Cost: 1.22},
			wantErr: false,
		},
		{
			name:    "new crypto create err",
			args:    args{shortTitle: "ETH", cost: -1.22},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCrypto(tt.args.shortTitle, tt.args.cost)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrypto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTitle(t *testing.T) {
	now := time.Now()
	cr := &Crypto{ShortTitle: "ETH", Cost: 1.22}
	cr.SetTitle("Ethereum")
	require.Equal(t, "Ethereum", cr.Title)
	cr.SetTimeStamp(now)
	require.Equal(t, now, cr.Created)
}
