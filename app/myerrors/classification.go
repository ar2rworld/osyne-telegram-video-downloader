package myerrors

type ErrorSeverity int

const (
	SeverityUser       ErrorSeverity = iota // Пользователь может исправить (плохой URL)
	SeverityMaintainer                      // Требуется действие админа (cookies)
	SeverityCritical                        // Сбой системы
)

type ClassifiedError interface {
	error
	Severity() ErrorSeverity
	UserMessage() string       // Безопасное сообщение для пользователя
	MaintainerMessage() string // Детальное сообщение для логов/алертов
}
