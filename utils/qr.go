package utils

import (
	"fmt"
	"image"
	"image/color"
)

/*
* @Author:hanyajun
* @Date:2019/4/30 14:56
* @Name:qrcode2console
* @Function: 二维码输出到console
* @SourceURL: https://github.com/Han-Ya-Jun/qrcode2console
 */

const (
	QR_CODE_SIZE        = 135 //原二维码图片的像素尺寸
	SHRINK_QR_CODE_SIZE = 39  //输出到控制台的二维码大小，单位为字符，包含两侧留白框
	MARGIN              = 12  //原二维码图片的边框像素尺寸
	MULTIPLE            = 3   //原二维码图片中每个小块的像素尺寸
)

type QRCode2Console struct {
	img             image.Image
	points          [QR_CODE_SIZE][QR_CODE_SIZE]int
	tmpShrinkPoints [QR_CODE_SIZE][SHRINK_QR_CODE_SIZE]int
	shrinkPoints    [SHRINK_QR_CODE_SIZE][SHRINK_QR_CODE_SIZE]int
}

// NewQRCode2ConsoleWithImage 通过二维码图片创建用于输出到控制台的二维码
func NewQRCode2ConsoleWithImage(imgin image.Image) *QRCode2Console {
	qr := &QRCode2Console{img: imgin}
	return qr
}

// binarization 二维码图片二值化 0－1
func (qr *QRCode2Console) binarization() {
	gray := image.NewGray(image.Rect(0, 0, QR_CODE_SIZE, QR_CODE_SIZE))
	for x := 0; x < QR_CODE_SIZE; x++ {
		for y := 0; y < QR_CODE_SIZE; y++ {
			r32, g32, b32, _ := qr.img.At(x, y).RGBA()
			r, g, b := int(r32>>8), int(g32>>8), int(b32>>8)
			if (r+g+b)/3 > 180 {
				qr.points[y][x] = 0
				gray.Set(x, y, color.Gray{uint8(255)})
			} else {
				qr.points[y][x] = 1
				gray.Set(x, y, color.Gray{uint8(0)})
			}
		}
	}
}

// shrink 缩小二值化数组
func (qr *QRCode2Console) shrink() {
	for x := 0; x < QR_CODE_SIZE; x++ {
		cal := 1 //不从0开始，留白
		for y := MARGIN + 1; y < QR_CODE_SIZE-MARGIN; y += MULTIPLE {
			qr.tmpShrinkPoints[x][cal] = qr.points[x][y]
			cal++
		}
	}
	for y := 1; y < SHRINK_QR_CODE_SIZE-1; y++ {
		row := 1
		for x := MARGIN + 1; x < QR_CODE_SIZE-MARGIN; x += MULTIPLE {
			qr.shrinkPoints[row][y] = qr.tmpShrinkPoints[x][y]
			row++
		}
	}
}

// Output 控制台输出二维码
func (qr *QRCode2Console) Output() {
	qr.binarization()
	qr.shrink()
	for x := 0; x < SHRINK_QR_CODE_SIZE; x++ {
		for y := 0; y < SHRINK_QR_CODE_SIZE; y++ {
			if qr.shrinkPoints[x][y] == 1 {
				fmt.Print("\033[40;40m  \033[0m")
				//randColor()
			} else {
				fmt.Print("\033[47;30m  \033[0m")
			}
		}
		fmt.Println()
	}
}
