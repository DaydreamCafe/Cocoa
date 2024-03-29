// Package charreverser 英文字符翻转
package charreverser

import (
	"regexp"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// commandRegex 命令正则表达式
const commandRegex = `[A-z]{1}([A-z]|\s)+[A-z]{1}`

var (
	charMap = map[rune]rune{
		'a': 'ɐ',
		'b': 'q',
		'c': 'ɔ',
		'd': 'p',
		'e': 'ǝ',
		'f': 'ɟ',
		'g': 'ƃ',
		'h': 'ɥ',
		'i': 'ᴉ',
		'j': 'ɾ',
		'k': 'ʞ',
		'l': 'l',
		'm': 'ɯ',
		'n': 'u',
		'o': 'o',
		'p': 'd',
		'q': 'b',
		'r': 'ɹ',
		's': 's',
		't': 'ʇ',
		'u': 'n',
		'v': 'ʌ',
		'w': 'ʍ',
		'x': 'x',
		'y': 'ʎ',
		'z': 'z',
		'A': '∀',
		'B': 'ᗺ',
		'C': 'Ɔ',
		'D': 'ᗡ',
		'E': 'Ǝ',
		'F': 'Ⅎ',
		'G': '⅁',
		'H': 'H',
		'I': 'I',
		'J': 'ſ',
		'K': 'ʞ',
		'L': '˥',
		'M': 'W',
		'N': 'N',
		'O': 'O',
		'P': 'Ԁ',
		'Q': 'Ò',
		'R': 'ᴚ',
		'S': 'S',
		'T': '⏊',
		'U': '∩',
		'V': 'Λ',
		'W': 'M',
		'X': 'X',
		'Y': '⅄',
		'Z': 'Z',
		' ': ' ',
	}

	// compiledRegex 编译后的正则表达式
	compiledRegex *regexp.Regexp = regexp.MustCompile(commandRegex)
)

// handleReverse 翻转命令handler
func handleReverse(ctx *zero.Ctx) {
	// 获取需要翻转的字符串
	results := compiledRegex.FindAllString(ctx.MessageString(), -1)

	str := results[0]

	// 将字符顺序翻转
	var tempStrBuilder strings.Builder
	for i := len(str) - 1; i >= 0; i-- {
		tempStrBuilder.WriteByte(str[i])
	}

	// 翻转字符字形
	var reversedStrBuilder strings.Builder
	for _, char := range tempStrBuilder.String() {
		reversedStrBuilder.WriteRune(charMap[char])
	}

	// 发送翻转后的字符串
	ctx.SendChain(message.Text(reversedStrBuilder.String()))
}
