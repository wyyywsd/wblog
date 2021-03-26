package helpers

import (
	"fmt"
	"gorm_demo/internal/models"
	"regexp"
	"strings"
)

func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func ReplaceHtml(articleContent string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	articleContent = re.ReplaceAllStringFunc(articleContent, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	articleContent = re.ReplaceAllString(articleContent, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	articleContent = re.ReplaceAllString(articleContent, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	articleContent = re.ReplaceAllString(articleContent, " ")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	articleContent = re.ReplaceAllString(articleContent, " ")
	return articleContent

}

func GetUserByComment(userId uint) models.User {
	user := models.FindUserById(fmt.Sprint(userId))
	return user
}
