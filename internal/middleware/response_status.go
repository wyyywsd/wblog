package middleware

/**
封装服务器返回状态吗
*/

//
//
//type Response struct {
//	Code int `json:"code"`
//	Msg string `json:"msg"`
//	Data interface{} `json:"data"`
//	Location string `json:"location"`
//}
//
//
//
////错误状态处理
//func ResposeError(context *gin.Context,code int,err error,msg string){
//	resp := &Response{Code:code,Msg:msg,Data:nil}
//	context.JSON(code,resp)
//	response,_:=json.Marshal(resp)
//	context.Set("response",string(response))
//	context.AbortWithError(code,err)
//}
//
////正确状态处理
//func ResponseSuccess(context *gin.Context,msg string,data interface{},location string) {
//	//resp := &Response{Code: 200, Msg: msg, Data: data,Location: location}
//	//context.JSON(200, resp)
//	//response, _ := json.Marshal(resp)
//	//context.Redirect(resp.Code, resp.Location)
//	context.Redirect(http.StatusMovedPermanently, "/index")
//}
