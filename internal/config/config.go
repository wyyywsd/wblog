package config

//type MyApiConfig struct {
//	Listen string `toml:"listen"`
//	DBServers map[string]db.DBServer `toml:"dbservers"`
//	RedisServers map[string]db.RedisServer `toml:"redisservers"`
//	UserAPI string `toml:"user_api"`
//
//}
////解析toml
//func UnmarshalConfig(tomlfile string) (*MyApiConfig,error) {
//	c:= &MyApiConfig{}
//	_,err := toml.DecodeFile(tomlfile,c)
//	if err != nil {
//		//解析错误
//		return c,err
//	}
//	return c,nil
//}
//
////获取psql数据库配置
//func(c MyApiConfig) DBServerConf(key string)(db.DBServer,bool){
//	s,ok := c.DBServers[key]
//	return s,ok
//}
//
////获取redis数据库的配置
//func (c MyApiConfig) RedisServerConf(key string)(db.RedisServer,bool){
//	s,ok := c.RedisServers[key]
//	return s,ok
//}
//
////监听地址
//func (c MyApiConfig) GetListenAddr() string {
//	return c.Listen
//}
//
////Validate 验证配置
//func (c *MyApiConfig) Validate() error {
//	if c.Listen == "" {
//		return fmt.Errorf("listen未配置")
//	}
//	if c.UserAPI == "" {
//		return fmt.Errorf("user_api未配置")
//	}
//	return 	nil
//}