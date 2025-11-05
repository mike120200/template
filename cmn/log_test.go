package cmn

import "testing"

func TestLoggerInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
	}
	if err := ViperInit(".conf_linux.json"); err != nil {
		t.Error(err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoggerInit(); (err != nil) != tt.wantErr {
				t.Errorf("LoggerInit() error = %v, wantErr %v", err, tt.wantErr)
			}

			Logger().Info("test")
			Logger().Error("test")
			Logger().Warn("test")
			Logger().Debug("test")
		})
	}
}
