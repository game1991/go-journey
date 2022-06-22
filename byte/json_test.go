package byte

import "testing"

func TestMashal(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "测试空对象",
			args:    args{b: []byte("")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Mashal(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Mashal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
