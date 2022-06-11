package types

type Logger interface {
	// AddBreadcrumb adds a breadcrumb which the logger can use to display additional log entry information
	AddBreadcrumb(breadcrumb string)
	// PopBreadcrumb pops the last breadcrumb ONLY if it matches the provided breadcrumb string
	PopBreadcrumb(breadcrumb string)
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(err error)
}
