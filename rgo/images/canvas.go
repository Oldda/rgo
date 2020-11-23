package images

import (
	"errors"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
)

//公共处理类
type Canvas struct {
	Pics []*Picture //画布内的图片元素
	Text []*Text    //画布内的文字元素

}

//实例化处理类
func NewCanvas() *Canvas {
	return &Canvas{
		Pics: make([]*Picture, 0),
		Text: make([]*Text, 0),
	}
}

//给画布追加图片元素
func (c *Canvas) AddImageElement(fileName string, background bool) *Picture {
	img := NewPicture(fileName, background)
	c.Pics = append(c.Pics, img)
	return img
}

//给画布追加文字元素
func (c *Canvas) AddTextElement(content string) *Text {
	text := NewText(content)
	c.Text = append(c.Text, text)
	return text
}

//水印
func (c *Canvas) Watermark(to string) {
	//验证是否有背景图
	var background image.Image
	elementPics := make([]*Picture, 0)
	for _, img := range c.Pics {
		if img.background {
			background = img.img
			continue
		}
		elementPics = append(elementPics, img)
	}
	if background == nil {
		handleErr(errors.New("无背景图"))
		return
	}
	//声明一个和背景图一样大小的内存空图
	img := image.NewRGBA(background.Bounds())

	//填充背景图
	draw.Draw(img, background.Bounds(), background, image.ZP, draw.Src)

	//合成除背景图之外的其他图片
	if len(elementPics) > 0 {
		var p_dx, p_dy int
		for _, p := range elementPics {
			//偏移量
			if p.Dx < 0 { //从右下角开始偏移
				p_dx = background.Bounds().Dx() + p.Dx
			} else {
				p_dx = p.Dx
			}
			if p.Dy < 0 {
				p_dy = background.Bounds().Dy() + p.Dy
			} else {
				p_dy = p.Dy
			}
			draw.Draw(img, background.Bounds().Add(image.Pt(p_dx, p_dy)), p.img, image.ZP, draw.Over)
		}
	}

	//合成文字
	if len(c.Text) > 0 {
		var t_dx, t_dy int
		for _, t := range c.Text {
			f := freetype.NewContext()
			f.SetDPI(108)
			f.SetFont(t.Font)
			f.SetFontSize(t.Size)
			f.SetClip(img.Bounds())
			f.SetDst(img)
			f.SetSrc(image.NewUniform(color.RGBA{R: t.R, G: t.G, B: t.B, A: t.A}))
			if t.Dx < 0 {
				t_dx = img.Bounds().Dx() + t.Dx
			} else {
				t_dx = t.Dx
			}
			if t.Dy < 0 {
				t_dy = img.Bounds().Dy() + t.Dy
			} else {
				t_dy = t.Dy
			}
			_, err := f.DrawString(t.Content, freetype.Pt(t_dx, t_dy))
			if err != nil {
				handleErr(err)
			}
		}
	}
	//存储图像
	save(to, img)
}

//合成-多张图合成一张图
func (c *Canvas) Compose(to string) {
	if len(c.Pics) <= 0 {
		handleErr(errors.New("最少传一张图片"))
	}
	//计算画布的大小 - 最大的坐标加上该图的长宽
	bigest_x := 0
	bigest_y := 0
	for _, pic := range c.Pics {
		if pic.Dx+pic.img.Bounds().Max.X > bigest_x {
			bigest_x = pic.Dx + pic.img.Bounds().Max.X
		}
		if pic.Dy+pic.img.Bounds().Max.Y > bigest_y {
			bigest_y = pic.Dy + pic.img.Bounds().Max.Y
		}
	}

	//生成该画布
	img := image.NewRGBA(image.Rect(0, 0, bigest_x, bigest_y))

	//根据每张图的坐标画图-拼接
	for _, pic := range c.Pics {
		draw.Draw(img, pic.img.Bounds().Add(image.Pt(pic.Dx, pic.Dy)), pic.img, image.ZP, draw.Over)
	}
	save(to, img)
}

//裁剪-单图
func (c *Canvas) Clip(to string, x0, y0, x1, y1 int) {
	pic := c.Pics[0]
	img := pic.clip(x0, y0, x1, y1)
	save(to, img)
}

//缩放-单图
func (c *Canvas) Scale(to string, width, height int) {
	pic := c.Pics[0]
	img := pic.scale(width, height)
	save(to, img)
}

//灰度
func (c *Canvas) ToGray(to string) {
	pic := c.Pics[0]
	img := pic.toGray()
	save(to, img)
}
