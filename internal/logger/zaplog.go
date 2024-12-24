package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLog обслуживает логгер на основе [Zap].
//
// [Zap]: https://github.com/uber-go/zap
type ZapLog struct {
	logger *zap.Logger
}

// NewZapLog создает логгер.
func NewZapLog(level string) (*ZapLog, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLog{
		logger: logger,
	}, nil
}

// Close очищает буфер, если он есть. Вызываем при завершении приложения.
func (l *ZapLog) Close() error {
	if err := l.logger.Sync(); err != nil {
		return err
	}

	return nil
}

// WithField добавляет поле в лог с указанным значением.
func (l *ZapLog) WithField(key string, value interface{}) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

// WithError добавляет ошибку в лог.
func (l *ZapLog) WithError(err error) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any("error", err.Error())),
	}
}

// Error пишет лог уровня error.
func (l *ZapLog) Error(msg string) {
	l.logger.Error(msg)
}

// Info пишет лог уровня info.
func (l *ZapLog) Info(msg string) {
	l.logger.Info(msg)
}

// Debug пишет лог уровня debug.
func (l *ZapLog) Debug(msg string) {
	l.logger.Debug(msg)
}
