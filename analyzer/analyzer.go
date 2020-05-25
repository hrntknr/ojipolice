package analyzer

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"

	"github.com/ikawaha/kagome.ipadic/splitter"
	"github.com/kyokomi/emoji"
)

type OjiLevel uint

const (
	Warn  = iota
	Alert = iota
	Safe  = iota
)

var emojiMap = emoji.CodeMap()

func em(str string) string {
	emoji, ok := emojiMap[str]
	if !ok {
		panic(fmt.Sprintf("emoji not found: %s", str))
	}
	return emoji
}

// パターンはおじさんの生態を完全に理解しているgithub.com/greymd/ojichatを参考

// こういうのこそ機械学習だって？？おじさんにそんなものは扱えないから知るか

// １人称は 名詞 代名詞 一般
// 半角カナつかったりしたらおじさん度は高め？？
var firstPerson = map[string]int{
	"ﾎﾞｸ":   3,
	"ｵﾚ":    3,
	"小生":    2,
	"オジサン":  3,
	"ｵｼﾞｻﾝ": 3,
	"おじさん":  3,
	"オイラ":   2,
}

// ﾁｬﾝはギルティ
var nameSuffix = map[string]int{
	"チャン": 3,
	"ﾁｬﾝ": 3,
	"ちゃん": 1,
}

// ココらへんは黒に限りなく近いグレー
var nanchatte = map[string]int{
	"ﾅﾝﾁｬｯﾃ": 3,
	"ナンチャッテ": 3,
	"なんちゃって": 3,
	"なんてね":   3,
	"冗談":     1,
}

// おじさんの聖地、これだけでは決定力にかける
var hotel = map[string]int{
	"ホテル": 2,
	"旅館":  2,
}

var date = map[string]int{
	"デート":  2,
	"カラオケ": 2,
	"ドライブ": 2,
}

var metaphor = map[string]int{
	"天使":   1,
	"女神":   1,
	"女優さん": 1,
	"お姫様":  1,
}

// おじさんは絵文字を連打するから割と小さめに設定するよ
// 上で除外した絵文字もいれるヨ！
// 開発OSで絵文字の処理依存するの草、おじさん多彩すぎだろ
var emojiList = map[string]int{
	// OTHER
	em(":hotel:"):      2,
	em(":love_hotel:"): 2,
	em(":microphone:"): 1,
	em(":blue_car:"):   1,
	em(":red_car:"):    1,
	// EMOJI_POS
	em(":smiley:"):              1,
	em(":raised_hand:"):         1,
	em(":exclamation:"):         3,
	em(":smile:"):               1,
	em(":laughing:"):            1,
	em(":kissing_closed_eyes:"): 1,
	em(":kissing_heart:"):       1,
	em(":two_hearts:"):          2,
	em(":heartpulse:"):          2,
	em(":heart_eyes:"):          2,
	em(":grin:"):                1,
	em(":yum:"):                 1,
	em(":joy:"):                 1,
	em(":blush:"):               1,
	em(":musical_note:"):        1,
	// EMOJI_NEG
	em(":sweat_drops:"):           1,
	em(":broken_heart:"):          1,
	em(":scream:"):                1,
	em(":cold_sweat:"):            1,
	em(":sob:"):                   1,
	em(":sweat:"):                 1,
	em(":persevere:"):             1,
	em(":confounded:"):            1,
	em(":disappointed_relieved:"): 1,
	em(":cry:"):                   1,
	// EMOJI_NEUT
	em(":zzz:"):                    1,
	em(":sleeping:"):               1,
	em(":slight_smile:"):           1,
	em(":money_mouth:"):            1,
	em(":sleepy:"):                 1,
	em(":sleeping_accommodation:"): 1,
	em(":sunglasses:"):             1,
	em(":triumph:"):                1,
	em(":unamused:"):               1,
	em(":kissing_smiling_eyes:"):   1,
	em(":smirk:"):                  1,
	em(":flushed:"):                1,
	em(":relieved:"):               1,
	// EMOJI_ASK
	em(":question:"):                     3,
	em(":interrobang:"):                  3,
	em(":thinking:"):                     1,
	em(":stuck_out_tongue_winking_eye:"): 1,
}

func CheckOjiLevel(content string) []OjiResult {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(splitter.ScanSentences)
	result := []OjiResult{}
	for scanner.Scan() {
		sentence := scanner.Text()
		result = append(result, checkOjiLevelWithSentence(sentence))
	}
	return result
}

func checkOjiLevelWithSentence(sentence string) OjiResult {
	ojiScore := 0

	// 末尾のカタカナの数をチェック
	endKatakana := 0
	buf := []rune(sentence)
	for i := len(buf); i > 0; i-- {
		if unicode.In(buf[i-1], unicode.Hiragana) {
			break
		}
		if unicode.In(buf[i-1], unicode.Katakana) {
			endKatakana++
		}
	}
	if endKatakana > 0 {
		ojiScore += 3
	}

	for w, score := range firstPerson {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range nameSuffix {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range nanchatte {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range hotel {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range date {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range metaphor {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	for w, score := range emojiList {
		if strings.Contains(sentence, w) {
			ojiScore += score
		}
	}

	if ojiScore >= 8 {
		return OjiResult{
			Level:    Alert,
			Sentence: sentence,
		}
	} else if ojiScore >= 4 {
		return OjiResult{
			Level:    Warn,
			Sentence: sentence,
		}
	} else {
		return OjiResult{
			Level:    Safe,
			Sentence: sentence,
		}
	}
}

type OjiResult struct {
	Level    OjiLevel
	Sentence string
}
