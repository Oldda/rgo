package images

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
)

type Text struct {
	Font     *truetype.Font //字体
	Size     float64 //文字大小
	Content  string  //文字内容
	Dx       int     //文字x轴留白距离
	Dy       int     //文字y轴留白距离
	R        uint8   //文字颜色值RGBA中的R值
	G        uint8   //文字颜色值RGBA中的G值
	B        uint8   //文字颜色值RGBA中的B值
	A        uint8   //文字颜色值RGBA中的A值
}

func NewText(text string)*Text{
	fontBytes, err := ioutil.ReadFile("./images/fonts/simkai.ttf")
	handleErr(err)
	font, err := freetype.ParseFont(fontBytes)
	handleErr(err)
	return &Text{
		Font: font,
		Size: 12,
		Content: text,
		Dx: 50,
		Dy: 50,
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
}
