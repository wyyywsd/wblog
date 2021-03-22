package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"time"

)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

func SafeURL(x string) template.URL {

	return template.URL(x)
}

func Test(x string) []string{
	b := strings.Split(x, "")


	fmt.Println(len(b))

	return b
}