package helpers

import (
	"html/template"
	"time"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

func SafeURL(x string) template.URL {
	return template.URL(x)
}
