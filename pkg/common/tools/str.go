package tools

import (
	"unicode"
)

// 日语字符范围
var japanese = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x3040, 0x309F, 1}, // 平假名
		{0x30A0, 0x30FF, 1}, // 片假名
		{0xFF66, 0xFF9D, 1}, // 半角片假名
	},
	R32: []unicode.Range32{
		{0x1B000, 0x1B0FF, 1}, // Kana Supplement
		{0x1B100, 0x1B12F, 1}, // Kana Extended-A
	},
}

// 中文字符范围
var chinese = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x4E00, 0x9FFF, 1}, // 常用汉字
	},
	R32: []unicode.Range32{
		{0x20000, 0x2A6DF, 1}, // 扩展 B
		{0x2A700, 0x2B73F, 1}, // 扩展 C
		{0x2B740, 0x2B81F, 1}, // 扩展 D
		{0x2B820, 0x2CEAF, 1}, // 扩展 E
		{0x2CEB0, 0x2EBEF, 1}, // 扩展 F
	},
}

func IsJapanString(str string) bool {
	rs := []rune(str)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if unicode.Is(japanese, r) {
			return true
		}
	}
	return false
}

// IsContainsEnglish checks if the string contains at least one English letter.
func IsContainsEnglish(s string) bool {
	for _, r := range s {
		// Check if the character is an English letter using unicode.IsLetter
		if unicode.IsLetter(r) && (unicode.Is(unicode.Latin, r) || unicode.Is(unicode.Greek, r)) {
			return true
		}
	}
	return false
}
