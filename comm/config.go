package comm

import "github.com/jinzhu/configor"

//配置读取
//import configor "github.com/jinzhu/configor"

func ReadConfig(config interface{}, files ...string) error {

	return configor.Load(config, files...)
	// data, err := ioutil.ReadFile(files[0])
	// fmt.Println(string(data))
	// err = json.Unmarshal(data, config)
	// return err
}
