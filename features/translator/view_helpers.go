package translator

// flagToLang maps country flag emoji to DeepL target language codes.
var flagToLang = map[string]string{
	"🇯🇵": "JA",
	"🇺🇸": "EN",
	"🇬🇧": "EN",
	"🇫🇷": "FR",
	"🇩🇪": "DE",
	"🇪🇸": "ES",
	"🇮🇹": "IT",
	"🇵🇹": "PT",
	"🇧🇷": "PT",
	"🇷🇺": "RU",
	"🇰🇷": "KO",
	"🇨🇳": "ZH",
	"🇹🇼": "ZH",
	"🇳🇱": "NL",
	"🇵🇱": "PL",
	"🇸🇪": "SV",
	"🇩🇰": "DA",
	"🇫🇮": "FI",
	"🇹🇷": "TR",
	"🇮🇩": "ID",
	"🇺🇦": "UK",
}

var langNames = map[string]string{
	"JA": "日本語",
	"EN": "英語",
	"FR": "フランス語",
	"DE": "ドイツ語",
	"ES": "スペイン語",
	"IT": "イタリア語",
	"PT": "ポルトガル語",
	"RU": "ロシア語",
	"KO": "韓国語",
	"ZH": "中国語",
	"NL": "オランダ語",
	"PL": "ポーランド語",
	"SV": "スウェーデン語",
	"DA": "デンマーク語",
	"FI": "フィンランド語",
	"TR": "トルコ語",
	"ID": "インドネシア語",
	"UK": "ウクライナ語",
}

func langName(code string) string {
	if name, ok := langNames[code]; ok {
		return name
	}
	return code
}
