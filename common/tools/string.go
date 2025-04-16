package tools

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"math/rand/v2"
	"strings"
	"unicode"
)

// uuid去除-，并将字符随机一半的概率转化为大写
func UUID() string {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	b := make([]byte, len(id))
	for i, c := range id {
		if unicode.IsLetter(c) && rand.IntN(10)%2 == 0 {
			b[i] = byte(unicode.ToUpper(c))
			continue
		}
		b[i] = id[i]
	}
	return string(b)
}

type RandomStringMod int

// 生成x位随机字符串，通过 mods 组合大写，小写，数字，特殊符号
const (
	LowerMod   RandomStringMod = 0 // 2^0 小写字母
	UpperMod   RandomStringMod = 2 // 2^1 大写字母
	NumberMod  RandomStringMod = 4 // 2^2 数字
	SymbolsMod RandomStringMod = 8 // 2^3 特殊符号

	AnyMod = LowerMod | UpperMod | NumberMod | SymbolsMod // 全类型
)

func RandomString(length int, mods RandomStringMod) string {
	var buf bytes.Buffer
	if mods&LowerMod == LowerMod {
		buf.WriteString("abcdefghijklmnopqrstuvwxyz")
	}
	if mods&UpperMod == UpperMod {
		buf.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if mods&NumberMod == NumberMod {
		buf.WriteString("0123456789")
	}
	if mods&SymbolsMod == SymbolsMod {
		buf.WriteString("~`@#$%^&*()-_=+[{]}|;:',<.>/?")
	}
	charset := buf.String()
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset)-1)]
	}
	return string(b)
}

const (
	ConcatWithComma     = "," // 逗号分割
	ConcatWithSemicolon = ";" // 逗号分割
)

func ConcatWith(source interface{}, target bytes.Buffer, split string) bytes.Buffer {
	if target.Len() == 0 {
		target.WriteString(fmt.Sprint(source))
		return target
	}
	target.WriteString(fmt.Sprintf("%s%v", split, source))
	return target
}
