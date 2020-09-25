package helpers

import (
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


func ReplaceHtml(article_content string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	article_content = re.ReplaceAllStringFunc(article_content, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	article_content = re.ReplaceAllString(article_content, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	article_content = re.ReplaceAllString(article_content, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	article_content = re.ReplaceAllString(article_content, " ")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	article_content = re.ReplaceAllString(article_content, " ")

	return article_content


}

