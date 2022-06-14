package initial

import (
	"fmt"
	"os"

	"gitee.com/don178/m3u8/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// Encoder:编码器(如何写入日志)。
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// WriterSyncer ：指定日志将写到哪里去
	f, err := os.OpenFile("_cast.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := zapcore.AddSync(f)

	// Log Level：哪种级别的日志将被写入。
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)

	global.Log = zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	global.Slog = global.Log.Sugar()
}
