package Utility

import "io/ioutil"

func LoadFile(fileName string) (string, error){
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
