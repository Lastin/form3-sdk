package form3_sdk

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		config SessionCofig
	}
	tests := []struct {
		name string
		args args
		want *SdkClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
