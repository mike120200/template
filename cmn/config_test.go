package cmn

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViperInit(t *testing.T) {
	type args struct {
		configName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test",
			args:    args{configName: ".conf_linux.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ViperInit(tt.args.configName); (err != nil) != tt.wantErr {
				t.Errorf("ViperInit() error = %v, wantErr %v", err, tt.wantErr)
			}
			result := viper.GetString("appServe.name")
			t.Logf("result: %s", result)
			assert.Equal(t, "template", result)
		})
	}
}
