package tools

import "os"

// IsFileExist 判断文件是否存在
func IsFileExist(path string) (bool, error) {
	_,err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true,nil

	//fileInfo, err := os.Stat(path)
	//if os.IsNotExist(err) {
	//	return false, nil
	//}
	////我这里判断了如果是0也算不存在
	//if fileInfo.Size() == 0 {
	//	return false, nil
	//}
	//if err == nil {
	//	return true, nil
	//}
	//return false, err
}
