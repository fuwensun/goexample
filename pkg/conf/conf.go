package conf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//
func GetConf(path string, data interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("readfile error [%w]", err)
		return err

	}
	err = yaml.Unmarshal(buf, data)
	if err != nil {
		err = fmt.Errorf("yaml unmarshal error [%w]", err)
		return err
	}
	return nil
}

// func printWd() {
// 	pwd, _ := os.Getwd()
// 	fileInfoList, err := ioutil.ReadDir(pwd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(len(fileInfoList))
// 	for i := range fileInfoList {
// 		fmt.Println(fileInfoList[i].Name())
// 	}
// }