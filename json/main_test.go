package main

import "testing"

func TestJSONPrt(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "测试结果",
		},
	}
	for _, tt := range tests {
		for i := 0; i < 100; i++ {
			t.Run(tt.name, func(t *testing.T) {
				JSONPrt()
			})
		}

	}
}
