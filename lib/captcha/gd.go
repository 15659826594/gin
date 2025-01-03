package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"os"
)

func imagettftext(img *image.RGBA, size float64, angle float64, x int, y int, color image.Image, fontFilename *truetype.Font, text string) {
	label := freetype.NewContext()
	label.SetDPI(72)
	label.SetFont(fontFilename)
	label.SetFontSize(size)
	label.SetClip(img.Bounds())
	label.SetDst(img)
	label.SetSrc(color)

	//deg := angle * math.Pi / 180.0

	pt := freetype.Pt(x, y)

	_, _ = label.DrawString(text, pt)
}

func imagecolorallocate(R int, G int, B int) image.Image {
	return image.NewUniform(color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255})
}

func imagestring(img *image.RGBA, level int, x int, y int, str string, color image.Image) {
	label := freetype.NewContext()
	label.SetDPI(72)
	label.SetClip(img.Bounds())
	label.SetDst(img)
	label.SetSrc(color)
	label.SetFontSize(float64(level))

	fontBytes, _ := os.ReadFile(__DIR__ + "/assets/ttfs/1.ttf")
	freeFont, _ := freetype.ParseFont(fontBytes)
	label.SetFont(freeFont)

	pt := freetype.Pt(x, y)
	_, _ = label.DrawString(str, pt)
}

func getimagesize(img *image.RGBA) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func imagesetpixel(img *image.RGBA, x int, y int, color image.Image) {
	img.Set(x, y, color.At(0, 0))
}

// imagerotate — 用给定角度旋转图像
// 将 image 图像按指定 angle 角度旋转。
// 旋转的中心是图像的中心，旋转后的图像可能与原始图像具有不同的尺寸。
func imagerotate(img *image.RGBA, angle float64, backgroundColor image.Image) {

}
