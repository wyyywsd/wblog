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
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	context.HTML(200, "picture_recognition.html", gin.H{
		"currentUser": currentUser,
		"userSession": currentUserName,
	})
}

//图片识别
func SubmitPictureRecognition(context *gin.Context) {
	picture64 := context.PostForm("picture_64")
	//<p><img src="data:image/png;base64,
	//" style="max-width:100%;"><br></p>
	picture64Re1 := strings.Replace(picture64, "<p><img src=\"data:image/png;base64,", "", -1)
	picture64Re2 := strings.Replace(picture64Re1, "\" style=\"max-width:100%;\"><br></p>", "", -1)
	urlValues := url.Values{"image": {picture64Re2}}
	accurateBasicUrl := viper.GetString("accurate_basic.url")
	accessToken := viper.GetString("accurate_basic.access_token")
	resp, err := http.PostForm(accurateBasicUrl+accessToken,
		urlValues)
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
	mapBody := string(body)
	//转换成 map
	var tempMap map[string][]map[string]string
	json.Unmarshal([]byte(mapBody), &tempMap)
	//对返回值做特殊处理  全部整合到一行
	//{"log_id": 1840946480854490027, "words_result_num": 4, "words_result": [{"words": "89860"}, {"words": "45124"}, {"words": "19C07"}, {"words": "72356"}]}
	listTemp := tempMap["words_result"]
	//[{"words": "89860"}, {"words": "45124"}, {"words": "19C07"}, {"words": "72356"}]
	iccids := ""
	for _, words := range listTemp {
		temp := words["words"]
		iccids = iccids + temp
	}
	fmt.Println(iccids)
	context.String(http.StatusOK, iccids)
}

//显示解绑批次列表
func ShowUnbindBatchIndex(context *gin.Context) {
	pageString := context.Param("page")
	//将page转换成int
	page, _ := strconv.Atoi(pageString)
	//获取所有的运营商类别 用于显示在主页面
	carriers, _ := models.AllCarriers()
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	//获取所有的批次
	unbindBatches, _ := models.FindUnbindBatchByPage(batchCount, page)

	//获取一共有多少批次
	count := models.UnbindCount()
	//通过批次的数量 算出分页一共有多少页   如果有余数  就加一
	pageCount := count / batchCount
	if count%batchCount != 0 {
		pageCount = (count / batchCount) + 1
	}
	context.HTML(200, "unbind_batch.html", gin.H{
		"carriers":      carriers,
		"currentUser":   currentUser,
		"userSession":   currentUserName,
		"unbindBatches": unbindBatches,
		"pageCount":     pageCount,
		"currentPage":   page,
	})
}

//新建一个解绑批次
func CreateUnbindBatch(context *gin.Context) {
	carrierIdString := context.PostForm("carrier_name")
	status := context.PostForm("batch_status")
	fmt.Println(status)
	carrierId, _ := strconv.Atoi(carrierIdString)
	models.CreateUnbindBatch(uint(carrierId), status)
	context.Redirect(http.StatusMovedPermanently, "/batch/index/1")
}

//显示某个解绑批次的所有卡号
func ShowUnbindBatch(context *gin.Context) {
	unbindBatchIdString := context.Param("unbind_batch_id")
	//获取的id 转换成int
	unbindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	unbindBatch, _ := models.FindUnbindBatchById(uint(unbindBatchId))
	//常规获取session
	session := sessions.Default(context)
	currentUserName := session.Get("sessionId")
	currentUser, _, _ := models.FindUserByUserName(fmt.Sprint(currentUserName))
	carrier := unbindBatch.FindCarrierByUnbindBatch()
	simCards, _ := models.FindSimCardsByUnbindBatch(unbindBatch)
	fmt.Println(len(simCards))
	context.HTML(200, "show_unbind_batch_detail.html", gin.H{
		"unbindBatch": unbindBatch,
		"currentUser": currentUser,
		"userSession": currentUserName,
		"simCards":    simCards,
		"carrier":     carrier,
	})
}

func NewSimCard(context *gin.Context) {
	agentName := context.PostForm("agent_name")
	iccid := context.PostForm("iccid")
	msisdn := context.PostForm("msisdn")
	unbindBatchIdString := context.Param("unbind_batch_id")
	replaceReason := context.PostForm("replace_reason")
	equipmentPhoto := context.PostForm("image_base64")
	//原本是图片直接用base64存入数据库 ， 现改为地址 ， 减少数据库的压力
	_, fileNameStr := WriteFile("file", equipmentPhoto)
	//fmt.Println(equipmentPhoto)
	file_name := ""
	if fileNameStr != "" {
		file_name = "/file" + fileNameStr
	}
	ubindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	models.CreateSimCards(agentName, iccid, msisdn, uint(ubindBatchId), replaceReason, file_name)
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbindBatchIdString+"")
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
	fileName := "/unbind_picture/" + date + "/" + curFileStr + strconv.Itoa(n) + "." + fileType
	file := path + fileName
	fmt.Println(file)
	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err := ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		fmt.Println("============================================================")
		log.Println(err)
		fmt.Println("============================================================")
	}

	return false, fileName

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
	simCardIdString := context.Param("sim_card_id")
	unbindBatchIdString := context.Param("unbind_batch_id")
	//将string类型的sim_card_id 转换成int
	simCardId, _ := strconv.Atoi(simCardIdString)
	//删除卡方法
	models.DeleteSimCardById(uint(simCardId))
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbindBatchIdString+"")
}

func UpdateUnbindBatchStatus(context *gin.Context) {
	unbindBatchIdString := context.Param("unbind_batch_id")
	batchStatus := context.PostForm("batch_status")
	unbindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	models.UpdateUnbindBatchStatusById(uint(unbindBatchId), batchStatus)
	context.Redirect(http.StatusMovedPermanently, "/show_unbind_batch/"+unbindBatchIdString+"")
}

func DeleteUnbindBatch(context *gin.Context) {
	unbindBatchIdString := context.Param("unbind_batch_id")
	unbindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	models.DeleteUnbindBatchById(uint(unbindBatchId))
	context.Redirect(http.StatusMovedPermanently, "/batch/index/1")
}

func ExportDataExcel(context *gin.Context) {
	unbindBatchIdString := context.Param("unbind_batch_id")
	unbindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	fmt.Println("进来了")
	unbindBatch, _ := models.FindUnbindBatchById(uint(unbindBatchId))
	carrier := unbindBatch.FindCarrierByUnbindBatch()
	simCardsChan := make(chan *models.SimCard, 40)
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
	var listForMapInOrder []string
	//用做判断， 无锡移动目前 解绑需要提供设备照片 ，所以单独判断  如果是无锡移动导出卡号，就导出图片  如果后续有其他运营商需要提供设备照片， 再修改这里的判断逻辑
	if carrier.Name == "无锡移动" {
		titles = map[string]string{"A1": "iccid", "B1": "msisdn", "C1": "图片"}
		listForMapInOrder = append(listForMapInOrder, "iccid", "msisdn", "picture")
	} else {
		titles = map[string]string{"A1": "iccid", "B1": "msisdn"}
		listForMapInOrder = append(listForMapInOrder, "iccid", "msisdn")
	}

	for key, value := range titles {
		f.SetCellValue("Sheet1", key, value)
	}
	//调用方法 去查询对应的卡号数据  把文件写入data中
	for i := 0; i < 2; i++ {
		go QuerySimCardDataForExcel(simCardsChan, &wg, &lock, &data, carrier)
	}

	//取到simcard数据 传入信道中
	simCardsTemp, _ := models.FindSimCardsByUnbindBatch(unbindBatch)
	for _, sim_card := range simCardsTemp {
		simCardsChan <- sim_card
	}
	close(simCardsChan)
	wg.Wait()
	//调用写入表格的方法
	DoExcel(&listForMapInOrder, b, f, &data)
	//保存文件
	filename := unbindBatch.FindCarrierByUnbindBatch().Name + "解绑.xlsx"
	filepath := "/Users/mac/Desktop/解绑专用/" + filename
	if err := f.SaveAs(filepath); err != nil {
		println(err.Error() + "123123123")
	}
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	context.Writer.Header().Add("Content-Type", "application/msexcel")
	context.File(filepath)
}

func QuerySimCardDataForExcel(sim_cards chan *models.SimCard, wg *sync.WaitGroup, lock *sync.Mutex, data *[](map[string]string), carrier models.Carrier) {
	fmt.Println("进来了2")
	defer wg.Done()
	for simCard := range sim_cards {
		dataItem := map[string]string{}
		dataItem["iccid"] = simCard.Iccid
		dataItem["msisdn"] = simCard.Msisdn
		//用做判断， 无锡移动目前 解绑需要提供设备照片 ，所以单独判断  如果是无锡移动导出卡号，就导出图片  如果后续有其他运营商需要提供设备照片， 再修改这里的判断逻辑
		if carrier.Name == "无锡移动" {
			dataItem["picture"] = simCard.EquipmentPhoto
		}
		//将数据存入data中
		lock.Lock()
		*data = append(*data, dataItem)
		fmt.Printf("记录了%s次数据\n", "***")
		lock.Unlock()
	}
}

//表格写入方法
func DoExcel(list_for_map_in_order *[]string, b int, f *excelize.File, data *[](map[string]string)) {
	//这个参数 临时用来判断 前一个卡号是否有图片添加，  如果有  图片后面的卡号  往下移动几行 再写入
	picTemp := 0
	for _, dataa := range *data {
		a := 1
		for _, key := range *list_for_map_in_order {
			col, _ := excelize.ColumnNumberToName(a)
			if key == "picture" {
				fmt.Println("进入图片这里了")
				if dataa[key] == "" {
					//到这里证明是写入卡号了 设置成0
					picTemp = 0
					break
				}
				fmt.Println(dataa[key])
				reg1 := regexp.MustCompile(`file\/unbind_picture.*`)
				fileName := reg1.FindAllStringSubmatch(dataa[key], -1)
				fmt.Println(fileName[0])
				fmt.Println(b)
				if err := f.AddPicture("Sheet1", col+strconv.Itoa(b), fileName[0][0], `{
        			"x_scale": 0.1,
        			"y_scale": 0.1
    			}`); err != nil {
					fmt.Println(err)
				}
				//到了这里证明  是添加了图片的 设置值为1
				picTemp = 1
			} else {
				if picTemp == 1 {
					b += 5
				}
				f.SetCellValue("Sheet1", col+strconv.Itoa(b), dataa[key])
				//到这里证明是写入卡号了 设置成0
				picTemp = 0
			}
			a++
		}
		b++
		fmt.Printf("执行了%d次\n", b)
	}
}

func ExportDataTxt(context *gin.Context) {
	unbindBatchIdString := context.Param("unbind_batch_id")
	unbindBatchId, _ := strconv.Atoi(unbindBatchIdString)
	fmt.Println("进来了")
	unbindBatch, _ := models.FindUnbindBatchById(uint(unbindBatchId))
	filename := unbindBatch.FindCarrierByUnbindBatch().Name + "解绑.txt"
	filepath := "/Users/mac/Desktop/解绑专用/" + filename

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		sim_cards, _ := models.FindSimCardsByUnbindBatch(unbindBatch)
		for _, sim_card := range sim_cards {
			f.Write([]byte(sim_card.Msisdn + "\r\n"))
		}

	}
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	context.Writer.Header().Add("Content-Type", "application/txt")
	context.File(filepath)
}
