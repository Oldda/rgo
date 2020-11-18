package images

import (
	"github.com/nfnt/resize"
	"image"
	"image/draw"
)

type Picture struct {
	img image.Image
	background bool //是否是背景图
	Dx int //从x轴位置起画
	Dy int //从y轴位置画起
}

func NewPicture(fileName string,background bool)*Picture{
	img := open(fileName)
	return &Picture{
		img: img,
		background: background,
		Dx: 0,
		Dy: 0,
	}
}

//单图片裁剪
func(p *Picture)clip(x0,y0,x1,y1 int)image.Image{
	img := image.NewRGBA(p.img.Bounds())
	draw.Draw(img, p.img.Bounds(), p.img, image.ZP, draw.Over)
	return img.SubImage(image.Rect(x0,y0,x1,y1))
}

//缩放
func(p *Picture)scale(width,height int)image.Image{
	return resize.Thumbnail(uint(width), uint(height), p.img, resize.Lanczos3)
}

//色素控制
func(p *Picture)toGray()image.Image{
	src := image.NewGray(p.img.Bounds())
	for x:=0;x<=src.Bounds().Dx();x++{
		for y:=0;y<=src.Bounds().Dy();y++{
			src.Set(x,y,p.img.At(x,y))
		}
	}
	return src
}