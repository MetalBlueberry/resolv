package aabb_test

import (
	"reflect"
	"testing"

	. "github.com/SolarLune/resolv/resolv/aabb"
)

func TestMerge(t *testing.T) {
	type args struct {
		a AABB
		b AABB
	}
	tests := []struct {
		name string
		args args
		want AABB
	}{
		{
			name: "Merge overlapping",
			args: args{
				a: &AABBData{0, 0, 3, 3},
				b: &AABBData{2, 0, 5, 2},
			},
			want: &AABBData{0, 0, 5, 3},
		},
		{
			name: "Merge next to each other",
			args: args{
				a: &AABBData{0, 0, 1, 1},
				b: &AABBData{2, 0, 3, 1},
			},
			want: &AABBData{0, 0, 3, 1},
		},
		{
			name: "Merge inside",
			args: args{
				a: &AABBData{-2, -2, 2, 2},
				b: &AABBData{-1, -1, 1, 1},
			},
			want: &AABBData{-2, -2, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
