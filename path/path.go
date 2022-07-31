package path

import (
	"fmt"
	"os"
	"strings"
)

func CheckPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CheckAndMakeSymlink(srcPath, symlinkPath string) {
	fileInfo, err := os.Stat(symlinkPath)

	if err != nil {
		if os.IsNotExist(err) {
			lastIndexByte := strings.LastIndexByte(symlinkPath, '/')
			fatherPath := symlinkPath[0:lastIndexByte]
			_, err = os.Stat(fatherPath)
			if err != nil {
				if os.IsNotExist(err) {
					os.MkdirAll(fatherPath, os.ModePerm)
				}
				os.Symlink(srcPath, symlinkPath)
				return
			} else {
				fmt.Printf("fail to get file:%s err:%v", fatherPath, err)
				return
			}
		} else {
			fmt.Printf("fail to get file:%s err:%v", symlinkPath, err)
			return
		}
	}

	if isDir := fileInfo.IsDir(); isDir {
		os.Remove(symlinkPath)
		os.Symlink(srcPath, symlinkPath)
		return
	}

	fileInfo, _ = os.Lstat(symlinkPath)
	mode := string(fileInfo.Mode().String()[0])
	if mode != "L" {
		os.RemoveAll(symlinkPath)
		os.Symlink(srcPath, symlinkPath)
		return
	}

	realPath, _ := os.Readlink(symlinkPath)
	if realPath != srcPath {
		os.Remove(symlinkPath)
		os.Symlink(srcPath, symlinkPath)
		return
	}
	return
}

func MkdirPath(folderPath string) error {
	pathInfo, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folderPath, 0777)
			if err != nil {
				fmt.Printf("fail to mkdir this path %s error %v", folderPath, err)
				return fmt.Errorf("fail to mkdir this path %s error %v", folderPath, err)
			}
			pathInfo, err = os.Stat(folderPath)
			if err != nil {
				fmt.Printf("fail to get this path %s error %v", folderPath, err)
				return fmt.Errorf("fail to get this path %s error %v", folderPath, err)
			}
		} else {
			fmt.Printf("fail to get this path %s error %v", folderPath, err)
			return fmt.Errorf("fail to get this path %s error %v", folderPath, err)
		}
	}
	if !pathInfo.IsDir() {
		fmt.Printf("the path %s is not folder", folderPath)
		return fmt.Errorf("the path %s is not folder", folderPath)
	}

	return nil
}
