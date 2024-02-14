package validator

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type testStruct struct {
		Name string `validate:"required" json:"name"`
		Age  int    `validate:"required,min=10" json:"age"`
	}

	type args struct {
		s any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{testStruct{
				Name: "",
				Age:  0,
			}},
			wantErr: true,
		},
		{
			name: "",
			args: args{testStruct{
				Name: "asdffff",
				Age:  1000,
			}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
