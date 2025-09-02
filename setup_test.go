package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_adaptLogLevelFromString(t *testing.T) {
	type args struct {
		lvl string
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{name: "debug", args: args{lvl: "debug"}, want: DebugLevel},
		{name: "info", args: args{lvl: "info"}, want: InfoLevel},
		{name: "warning", args: args{lvl: "warning"}, want: WarnLevel},
		{name: "error", args: args{lvl: "error"}, want: ErrorLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, _ := adaptLogLevelFromString(tt.args.lvl)
			assert.Equalf(t, tt.want, level, "adaptLogLevelFromString(%v)", tt.args.lvl)
		})
	}
}
