package untils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var (
	Iv_aes  = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7}
	Key_aes = "ytUQ7l2ZZu8mLvJZ"
)

var (
	Iv_des  = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	Key_des = "b3L26XNL"
)

type JWData struct {
	Lng float64 // 经纬度
	Lat float64
}

type CtiyJWData struct {
	Location   JWData
	Precise    int
	Confidence int
	Level      string
}

type BodyData struct {
	Status string
	Result CtiyJWData
}

func Pksc7padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func Pksc7unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

func EncryptAES(src []byte) []byte {
	block, _ := aes.NewCipher([]byte(Key_aes))
	src = Pksc7padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, Iv_aes)
	blockmode.CryptBlocks(src, src)
	return src
}

func DecryptAES(src []byte) []byte {
	block, _ := aes.NewCipher([]byte(Key_aes))
	blockmode := cipher.NewCBCDecrypter(block, Iv_aes)
	blockmode.CryptBlocks(src, src)
	src = Pksc7unpadding(src)
	return src
}

func Pksc5padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func Pksc5unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func EncryptDES(src []byte) []byte {
	block, _ := des.NewCipher([]byte(Key_des))
	src = Pksc5padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, Iv_des)
	blockmode.CryptBlocks(src, src)
	return src
}

func DecryptDES(src []byte) []byte {
	block, _ := des.NewCipher([]byte(Key_des))
	blockmode := cipher.NewCBCDecrypter(block, Iv_des)
	blockmode.CryptBlocks(src, src)
	src = Pksc5unpadding(src)
	return src
}

func RandStrings(lenth int) string {
	chars := []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "v", "K", "M", "N", "P", "Q", "R", "S", "T", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "h", "i", "j", "k", "m", "n", "p", "r", "s", "t", "w", "x", "y", "z", "2", "3", "4", "5", "6", "7", "8"}
	restr := ""
	for i := 0; i < lenth; i++ {
		restr += chars[rand.Intn(len(chars))]
	}
	return restr
}

/**
生成16 md5
*/
func MD5_16(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))[8:24]
	//return hex.EncodeToString(ctx.Sum(nil))	32md5
}
func MD5_32(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	//return hex.EncodeToString(ctx.Sum(nil))[8:24]
	return hex.EncodeToString(ctx.Sum(nil))
}

/**
根据城市获取经纬度
*/
func GetCityToLocation(strCtiy string) (float64, float64) {
	resp, err := http.Get("http://api.map.baidu.com/geocoder?address=" + strCtiy + "&output=json&key=pckg0S4gcS65cSZbRdlxyb4kTq3DIAsQ&city=" + strCtiy)
	if err != nil {
		fmt.Println("http.Get err:", err.Error())
		return 0.0, 0.0
	}
	body, errbody := ioutil.ReadAll(resp.Body)
	if errbody != nil {
		fmt.Println("ioutil.ReadAll errbody:", errbody.Error())
		return 0.0, 0.0
	}
	//fmt.Println("body:", string(body))
	// 解析数据
	st := &BodyData{}
	json.Unmarshal(body, &st)
	return st.Result.Location.Lat, st.Result.Location.Lng
}
