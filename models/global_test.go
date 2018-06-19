package models

import (
	"reflect"
	"testing"
)

func TestGeneratePagination(t *testing.T) {
	type args struct {
		data  interface{}
		count int
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "test",
			args: args{
				data: []string{
					"first",
					"secend",
				},
				count: 2,
				page:  1,
				limit: 10,
			},
			want: Pagination{
				Data: []string{
					"first",
					"secend",
				},
				Pages: PaginationPages{
					Current: 1,
					Prev:    0,
					HasPrev: false,
					Next:    2,
					HasNext: false,
					Total:   1,
				},
				Items: PaginationItems{
					Limit: 10,
					Begin: 1,
					End:   2,
					Total: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePagination(tt.args.data, tt.args.count, tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}
