package comm

import (
	"os"
)

var UnityToolHelper unityToolHelper

type unityToolHelper struct{}

// func(unityToolHelper)RemoveFile(filepath string){
//     ioutil
// }
//FileExist 文件是否存在
func (unityToolHelper) FileExist(filepath string) (f os.FileInfo, exist bool) {
	exist = true
	var err error
	if f, err = os.Stat(filepath); os.IsNotExist(err) {
		exist = false
	}
	return f, exist
}

//MKdir 创建目录
func (unityToolHelper) MKdir(dir string) error {
	err := os.Mkdir(dir, os.ModePerm)
	return err
}
