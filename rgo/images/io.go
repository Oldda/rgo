package images

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

//打开文件并解析为内存图像
func open(fileName string)image.Image{
	file,err := os.Open(fileName)
	handleErr(err)
	defer file.Close()

	var img image.Image
	suffix := getImageSuffix(fileName)
	if suffix == "jpg" || suffix == "jpeg" {
		img, err = jpeg.Decode(file)
	} else if suffix == "png" {
		img, err = png.Decode(file)
	}
	handleErr(err)
	return img
}

//保存图像
func save(fileName string, img image.Image){
	newFile, err := os.Create(fileName)
	handleErr(err)
	defer newFile.Close()


	suffix := getImageSuffix(fileName)
	if suffix == "jpg" || suffix == "jpeg" {
		err = jpeg.Encode(newFile, img, nil)
	} else if suffix == "png" {
		err = png.Encode(newFile,img)
	}
	handleErr(err)
}

//获取图片的后缀
func getImageSuffix(fileName string)string{
	if strings.HasSuffix(fileName,"jpg"){
		return "jpg"
	}
	if strings.HasSuffix(fileName,"jpeg") {
		return "jpeg"
	}
	if strings.HasSuffix(fileName,"png"){
		return "png"
	}
	handleErr(errors.New("不支持的后缀类型"))
	return ""
}
