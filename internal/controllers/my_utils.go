package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm_demo/internal/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)
//跳转到图片识别界面
func PictureRecognition(context *gin.Context){
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	context.HTML(200, "picture_recognition.html", gin.H{
		"current_user":current_user,
		"user_session":current_user_name,
	})
}
//图片识别
func SubmitPictureRecognition(context *gin.Context){
	picture_64 := context.PostForm("picture_64")
	//<p><img src="data:image/png;base64,
	//" style="max-width:100%;"><br></p>
	picture_64_re1 := strings.Replace(picture_64,"<p><img src=\"data:image/png;base64,","",-1)
	picture_64_re2 := strings.Replace(picture_64_re1,"\" style=\"max-width:100%;\"><br></p>","",-1)
	url_values := url.Values{"image": {picture_64_re2}}
	resp, err := http.PostForm("https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic?access_token=24.85ab595b0905856f5f7b8552de2dfd38.2592000.1604914532.282335-20446596",
		url_values)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))
	map_body := string(body)

	//转换成 map
	var tempMap map[string][]map[string]string
	json.Unmarshal([]byte(map_body), &tempMap)

	//对返回值做特殊处理  全部整合到一行
	//{"log_id": 1840946480854490027, "words_result_num": 4, "words_result": [{"words": "89860"}, {"words": "45124"}, {"words": "19C07"}, {"words": "72356"}]}
	list_temp := tempMap["words_result"]
	//[{"words": "89860"}, {"words": "45124"}, {"words": "19C07"}, {"words": "72356"}]
	iccids := ""
	for _,words := range list_temp{
		temp := words["words"]
		iccids = iccids+temp
	}
fmt.Println(iccids)
	context.String(http.StatusOK,iccids)
}
