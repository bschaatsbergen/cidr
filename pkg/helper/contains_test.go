package helper_test

import (
	"testing"

	"github.com/bschaatsbergen/cidr/pkg/helper"
)

func TestContainsInt(t *testing.T) {
	type args struct {
		ints         []int
		specifiedInt int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ContainsInt() should return true",
			args: args{
				ints:         []int{1, 2, 3, 4, 5},
				specifiedInt: 3,
			},
			want: true,
		},
		{
			name: "ContainsInt() should return false",
			args: args{
				ints:         []int{1, 2, 3, 4, 5},
				specifiedInt: 6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.ContainsInt(tt.args.ints, tt.args.specifiedInt); got != tt.want {
				t.Errorf("ContainsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
