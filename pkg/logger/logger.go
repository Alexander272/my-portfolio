package logger

type Logger interface {
	Debug(msg string, params map[string]interface{})
	Info(msg string, params map[string]interface{})
	Error(msg string, params map[string]interface{})
	Fatal(msg string, params map[string]interface{})
}
