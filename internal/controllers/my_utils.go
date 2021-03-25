package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm_demo/internal/models"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//分页的数量
const batchCount int = 10

//跳转到图片识别界面
func PictureRecognition(context *gin.Context) {
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	context.HTML(200, "picture_recognition.html", gin.H{
		"current_user": current_user,
		"user_session": current_user_name,
	})
}

//图片识别
func SubmitPictureRecognition(context *gin.Context) {
	picture_64 := context.PostForm("picture_64")
	//<p><img src="data:image/png;base64,
	//" style="max-width:100%;"><br></p>
	picture_64_re1 := strings.Replace(picture_64, "<p><img src=\"data:image/png;base64,", "", -1)
	picture_64_re2 := strings.Replace(picture_64_re1, "\" style=\"max-width:100%;\"><br></p>", "", -1)
	url_values := url.Values{"image": {picture_64_re2}}
	accurate_basic_url := viper.GetString("accurate_basic.url")
	access_token := viper.GetString("accurate_basic.access_token")
	resp, err := http.PostForm(accurate_basic_url+access_token,
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
	for _, words := range list_temp {
		temp := words["words"]
		iccids = iccids + temp
	}
	fmt.Println(iccids)
	context.String(http.StatusOK, iccids)
}

//显示解绑批次列表
func ShowUnbindBatchIndex(context *gin.Context) {
	page_string := context.Param("page")
	//将page转换成int
	page, _ := strconv.Atoi(page_string)
	//获取所有的运营商类别 用于显示在主页面
	carriers, _ := models.AllCarriers()
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	//获取所有的批次
	unbind_batchs, _ := models.FindUnbindBatchByPage(batchCount, page)

	//获取一共有多少批次
	count := models.UnbindCount()
	//通过批次的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count / batchCount
	if count%batchCount != 0 {
		pageCount = (count / batchCount) + 1
	}
	context.HTML(200, "unbind_batch.html", gin.H{
		"carriers":      carriers,
		"current_user":  current_user,
		"user_session":  current_user_name,
		"unbind_batchs": unbind_batchs,
		"pageCount":     pageCount,
		"current_page":  page,
	})
}

//新建一个解绑批次
func CreateUnbindBatch(context *gin.Context) {
	carrier_id_string := context.PostForm("carrier_name")
	status := context.PostForm("batch_status")
	fmt.Println(status)
	//session := sessions.Default(context)
	//current_user_name := session.Get("sessionid")
	//current_user,_,_:= models.FindUserByUserName(fmt.Sprint(current_user_name))
	carrier_id, _ := strconv.Atoi(carrier_id_string)
	models.CreateUnbindBatch(uint(carrier_id), status)

	//carriers,_ := models.AllCarriers()

	context.Redirect(http.StatusMovedPermanently, "/batch/index/1")
}

//显示某个解绑批次的所有卡号
func ShowUnbindBatch(context *gin.Context) {
	unbind_batch_id_string := context.Param("unbind_batch_id")
	//获取的id 转换成int
	unbind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	unbind_batch, _ := models.FindUnbindBatchById(uint(unbind_batch_id))
	//常规获取session
	session := sessions.Default(context)
	current_user_name := session.Get("sessionid")
	current_user, _, _ := models.FindUserByUserName(fmt.Sprint(current_user_name))
	carrier := unbind_batch.FindCarrierByUnbindBatch()
	sim_cards, _ := models.FindSimCardsByUnbindBatch(unbind_batch)
	fmt.Println(len(sim_cards))
	context.HTML(200, "show_unbind_batch_detail.html", gin.H{
		"unbind_batch": unbind_batch,
		"current_user": current_user,
		"user_session": current_user_name,
		"sim_cards":    sim_cards,
		"carrier":      carrier,
	})
}

func NewSimCard(context *gin.Context) {
	agent_name := context.PostForm("agent_name")
	iccid := context.PostForm("iccid")
	msisdn := context.PostForm("msisdn")
	unbind_batch_id_string := context.Param("unbind_batch_id")
	replace_reason := context.PostForm("replace_reason")
	equipment_photo := context.PostForm("image_base64")
	//原本是图片直接用base64存入数据库 ， 现改为地址 ， 减少数据库的压力
	_, file_name_s := WriteFile("file", equipment_photo)
	//fmt.Println(equipment_photo)
	file_name := ""
	if file_name_s != "" {
		file_name = "/file" + file_name_s
	}
	ubind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	models.CreateSimCards(agent_name, iccid, msisdn, uint(ubind_batch_id), replace_reason, file_name)
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbind_batch_id_string+"")
}

//base64 图片解码存入服务器
func WriteFile(path string, base64_image_content string) (bool, string) {
	fmt.Println("第一步成功")
	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)
	if !b {
		return false, ""
	}
	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)

	allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	fmt.Println(allData)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取
	fmt.Println(fileType)
	base64Str := re.ReplaceAllString(base64_image_content, "")
	//fmt.Println(base64Str)
	date := time.Now().Format("2006-01-02")
	if ok := IsFileExist(path + "/" + date); !ok {
		os.Mkdir(path+"/unbind_picture/"+date, 0666)
	}
	curFileStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(99999)
	file_name := "/unbind_picture/" + date + "/" + curFileStr + strconv.Itoa(n) + "." + fileType
	file := path + file_name
	fmt.Println(file)
	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err := ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		fmt.Println("============================================================")
		log.Println(err)
		fmt.Println("============================================================")
	}

	return false, file_name

}

//判断文件是否存在

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

func DeleteSimCard(context *gin.Context) {
	sim_card_id_string := context.Param("sim_card_id")
	unbind_batch_id_string := context.Param("unbind_batch_id")
	//将string类型的sim_card_id 转换成int
	sim_card_id, _ := strconv.Atoi(sim_card_id_string)
	//删除卡方法
	models.DeleteSimCardById(uint(sim_card_id))
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbind_batch_id_string+"")
}

func UpdateUnbindBatchStatus(context *gin.Context) {
	unbind_batch_id_string := context.Param("unbind_batch_id")
	batch_status := context.PostForm("batch_status")
	unbind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	models.UpdateUnbindBatchStatusById(uint(unbind_batch_id), batch_status)
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbind_batch_id_string+"")
}

func DeleteUnbindBatch(context *gin.Context) {
	unbind_batch_id_string := context.Param("unbind_batch_id")
	unbind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	models.DeleteUnbindBatchById(uint(unbind_batch_id))
	context.Redirect(http.StatusMovedPermanently, "/batch/index/1")
}

func ExportDataExcel(context *gin.Context) {
	unbind_batch_id_string := context.Param("unbind_batch_id")
	unbind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	fmt.Println("进来了")
	unbind_batch, _ := models.FindUnbindBatchById(uint(unbind_batch_id))
	carrier := unbind_batch.FindCarrierByUnbindBatch()
	sim_cards_chan := make(chan *models.SimCard, 40)
	var wg sync.WaitGroup
	wg.Add(2)
	lock := sync.Mutex{}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	b := 2
	//创建一个切片 里面存的是map  用于把最终信息保存下来，遍历写到表格中
	var data [](map[string]string)
	//表格表头
	titles := map[string]string{}
	//创建一个切片 用于对无序的map进行有序的输出
	var list_for_map_in_order []string
	//用做判断， 无锡移动目前 解绑需要提供设备照片 ，所以单独判断  如果是无锡移动导出卡号，就导出图片  如果后续有其他运营商需要提供设备照片， 再修改这里的判断逻辑
	if carrier.Name == "无锡移动" {
		titles = map[string]string{"A1": "iccid", "B1": "msisdn", "C1": "图片"}
		list_for_map_in_order = append(list_for_map_in_order, "iccid", "msisdn", "picture")
	} else {
		titles = map[string]string{"A1": "iccid", "B1": "msisdn"}
		list_for_map_in_order = append(list_for_map_in_order, "iccid", "msisdn")
	}

	for key, value := range titles {
		f.SetCellValue("Sheet1", key, value)
	}
	//调用方法 去查询对应的卡号数据  把文件写入data中
	for i := 0; i < 2; i++ {
		go QuerySimCardData_Wangzq(sim_cards_chan, &wg, &lock, &data, carrier)
	}

	//取到simcard数据 传入信道中
	sim_cardss, _ := models.FindSimCardsByUnbindBatch(unbind_batch)
	for _, sim_card := range sim_cardss {
		sim_cards_chan <- sim_card
	}
	close(sim_cards_chan)
	wg.Wait()
	//调用写入表格的方法
	DoExcel(&list_for_map_in_order, b, f, &data)
	//保存文件
	filename := unbind_batch.FindCarrierByUnbindBatch().Name + "解绑.xlsx"
	filepath := "/Users/mac/Desktop/解绑专用/" + filename
	if err := f.SaveAs(filepath); err != nil {
		println(err.Error() + "123123123")
	}
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	context.Writer.Header().Add("Content-Type", "application/msexcel")
	context.File(filepath)
}

func QuerySimCardData_Wangzq(sim_cards chan *models.SimCard, wg *sync.WaitGroup, lock *sync.Mutex, data *[](map[string]string), carrier models.Carrier) {
	fmt.Println("进来了2")
	defer wg.Done()
	for sim_card := range sim_cards {
		data_item := map[string]string{}
		data_item["iccid"] = sim_card.Iccid
		data_item["msisdn"] = sim_card.Msisdn
		//用做判断， 无锡移动目前 解绑需要提供设备照片 ，所以单独判断  如果是无锡移动导出卡号，就导出图片  如果后续有其他运营商需要提供设备照片， 再修改这里的判断逻辑
		if carrier.Name == "无锡移动" {
			data_item["picture"] = sim_card.EquipmentPhoto
		}
		//将数据存入data中
		lock.Lock()
		*data = append(*data, data_item)
		fmt.Printf("记录了%s次数据\n", "***")
		lock.Unlock()
	}
}

//表格写入方法
func DoExcel(list_for_map_in_order *[]string, b int, f *excelize.File, data *[](map[string]string)) {
	//这个参数 临时用来判断 前一个卡号是否有图片添加，  如果有  图片后面的卡号  往下移动几行 再写入
	pic_temp := 0
	for _, dataa := range *data {
		a := 1
		for _, key := range *list_for_map_in_order {
			col, _ := excelize.ColumnNumberToName(a)
			if key == "picture" {
				fmt.Println("进入图片这里了")
				if dataa[key] == "" {
					//到这里证明是写入卡号了 设置成0
					pic_temp = 0
					break
				}
				fmt.Println(dataa[key])
				reg1 := regexp.MustCompile(`file\/unbind_picture.*`)
				file_name := reg1.FindAllStringSubmatch(dataa[key], -1)
				fmt.Println(file_name[0])
				fmt.Println(b)
				if err := f.AddPicture("Sheet1", col+strconv.Itoa(b), file_name[0][0], `{
        			"x_scale": 0.1,
        			"y_scale": 0.1
    			}`); err != nil {
					fmt.Println(err)
				}
				//到了这里证明  是添加了图片的 设置值为1
				pic_temp = 1
			} else {
				if pic_temp == 1 {
					b += 5
				}
				f.SetCellValue("Sheet1", col+strconv.Itoa(b), dataa[key])
				//到这里证明是写入卡号了 设置成0
				pic_temp = 0
			}
			a++
		}
		b++
		fmt.Printf("执行了%d次\n", b)
	}
}

func ExportDataTxt(context *gin.Context) {
	unbind_batch_id_string := context.Param("unbind_batch_id")
	unbind_batch_id, _ := strconv.Atoi(unbind_batch_id_string)
	fmt.Println("进来了")
	unbind_batch, _ := models.FindUnbindBatchById(uint(unbind_batch_id))
	filename := unbind_batch.FindCarrierByUnbindBatch().Name + "解绑.txt"
	filepath := "/Users/mac/Desktop/解绑专用/" + filename

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		sim_cards, _ := models.FindSimCardsByUnbindBatch(unbind_batch)
		for _, sim_card := range sim_cards {
			f.Write([]byte(sim_card.Msisdn + "\r\n"))
		}

	}
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	context.Writer.Header().Add("Content-Type", "application/txt")
	context.File(filepath)
}
