// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

// FlagToLang maps country flag emoji to DeepL target language codes.
var FlagToLang = map[string]string{
	"\U0001F1E6\U0001F1E9": "CA",  // 🇦🇩 Andorra → カタルーニャ語
	"\U0001F1E6\U0001F1EB": "PRS", // 🇦🇫 Afghanistan → ダリー語
	"\U0001F1E6\U0001F1F1": "SQ",  // 🇦🇱 Albania → アルバニア語
	"\U0001F1E6\U0001F1F2": "HY",  // 🇦🇲 Armenia → アルメニア語
	"\U0001F1E6\U0001F1FF": "AZ",  // 🇦🇿 Azerbaijan → アゼルバイジャン語
	"\U0001F1E7\U0001F1E6": "BS",  // 🇧🇦 Bosnia → ボスニア語
	"\U0001F1E7\U0001F1E9": "BN",  // 🇧🇩 Bangladesh → ベンガル語
	"\U0001F1E7\U0001F1EC": "BG",  // 🇧🇬 Bulgaria → ブルガリア語
	"\U0001F1E7\U0001F1F7": "PT",  // 🇧🇷 Brazil → ポルトガル語
	"\U0001F1E7\U0001F1FC": "TN",  // 🇧🇼 Botswana → ツワナ語
	"\U0001F1E7\U0001F1FE": "BE",  // 🇧🇾 Belarus → ベラルーシ語
	"\U0001F1E8\U0001F1E9": "LN",  // 🇨🇩 DRC → リンガラ語
	"\U0001F1E8\U0001F1F3": "ZH",  // 🇨🇳 China → 中国語
	"\U0001F1E8\U0001F1FF": "CS",  // 🇨🇿 Czech Republic → チェコ語
	"\U0001F1E9\U0001F1EA": "DE",  // 🇩🇪 Germany → ドイツ語
	"\U0001F1E9\U0001F1F0": "DA",  // 🇩🇰 Denmark → デンマーク語
	"\U0001F1EA\U0001F1EA": "ET",  // 🇪🇪 Estonia → エストニア語
	"\U0001F1EA\U0001F1F8": "ES",  // 🇪🇸 Spain → スペイン語
	"\U0001F1EB\U0001F1EE": "FI",  // 🇫🇮 Finland → フィンランド語
	"\U0001F1EB\U0001F1F7": "FR",  // 🇫🇷 France → フランス語
	"\U0001F1EC\U0001F1E7": "EN",  // 🇬🇧 UK → 英語
	"\U0001F1EC\U0001F1EA": "KA",  // 🇬🇪 Georgia → ジョージア語
	"\U0001F1EC\U0001F1F7": "EL",  // 🇬🇷 Greece → ギリシャ語
	"\U0001F1ED\U0001F1F0": "YUE", // 🇭🇰 Hong Kong → 広東語
	"\U0001F1ED\U0001F1F7": "HR",  // 🇭🇷 Croatia → クロアチア語
	"\U0001F1ED\U0001F1F9": "HT",  // 🇭🇹 Haiti → ハイチ・クレオール語
	"\U0001F1ED\U0001F1FA": "HU",  // 🇭🇺 Hungary → ハンガリー語
	"\U0001F1EE\U0001F1E9": "ID",  // 🇮🇩 Indonesia → インドネシア語
	"\U0001F1EE\U0001F1EA": "GA",  // 🇮🇪 Ireland → アイルランド語
	"\U0001F1EE\U0001F1F1": "HE",  // 🇮🇱 Israel → ヘブライ語
	"\U0001F1EE\U0001F1F3": "HI",  // 🇮🇳 India → ヒンディー語
	"\U0001F1EE\U0001F1F7": "FA",  // 🇮🇷 Iran → ペルシア語
	"\U0001F1EE\U0001F1F8": "IS",  // 🇮🇸 Iceland → アイスランド語
	"\U0001F1EE\U0001F1F9": "IT",  // 🇮🇹 Italy → イタリア語
	"\U0001F1EF\U0001F1F5": "JA",  // 🇯🇵 Japan → 日本語
	"\U0001F1F0\U0001F1EC": "KY",  // 🇰🇬 Kyrgyzstan → キルギス語
	"\U0001F1F0\U0001F1F7": "KO",  // 🇰🇷 South Korea → 韓国語
	"\U0001F1F0\U0001F1FF": "KK",  // 🇰🇿 Kazakhstan → カザフ語
	"\U0001F1F1\U0001F1F8": "ST",  // 🇱🇸 Lesotho → ソト語
	"\U0001F1F1\U0001F1F9": "LT",  // 🇱🇹 Lithuania → リトアニア語
	"\U0001F1F1\U0001F1FA": "LB",  // 🇱🇺 Luxembourg → ルクセンブルク語
	"\U0001F1F1\U0001F1FB": "LV",  // 🇱🇻 Latvia → ラトビア語
	"\U0001F1F2\U0001F1EC": "MG",  // 🇲🇬 Madagascar → マダガスカル語
	"\U0001F1F2\U0001F1F0": "MK",  // 🇲🇰 North Macedonia → マケドニア語
	"\U0001F1F2\U0001F1F2": "MY",  // 🇲🇲 Myanmar → ビルマ語
	"\U0001F1F2\U0001F1F3": "MN",  // 🇲🇳 Mongolia → モンゴル語
	"\U0001F1F2\U0001F1F9": "MT",  // 🇲🇹 Malta → マルタ語
	"\U0001F1F2\U0001F1FE": "MS",  // 🇲🇾 Malaysia → マレー語
	"\U0001F1F3\U0001F1F1": "NL",  // 🇳🇱 Netherlands → オランダ語
	"\U0001F1F3\U0001F1F4": "NO",  // 🇳🇴 Norway → ノルウェー語
	"\U0001F1F3\U0001F1F5": "NE",  // 🇳🇵 Nepal → ネパール語
	"\U0001F1F3\U0001F1FF": "MI",  // 🇳🇿 New Zealand → マオリ語
	"\U0001F1F5\U0001F1EA": "QU",  // 🇵🇪 Peru → ケチュア語
	"\U0001F1F5\U0001F1ED": "TL",  // 🇵🇭 Philippines → タガログ語
	"\U0001F1F5\U0001F1F0": "UR",  // 🇵🇰 Pakistan → ウルドゥー語
	"\U0001F1F5\U0001F1F1": "PL",  // 🇵🇱 Poland → ポーランド語
	"\U0001F1F5\U0001F1F9": "PT",  // 🇵🇹 Portugal → ポルトガル語
	"\U0001F1F5\U0001F1FE": "GN",  // 🇵🇾 Paraguay → グアラニー語
	"\U0001F1F7\U0001F1F4": "RO",  // 🇷🇴 Romania → ルーマニア語
	"\U0001F1F7\U0001F1F8": "SR",  // 🇷🇸 Serbia → セルビア語
	"\U0001F1F7\U0001F1FA": "RU",  // 🇷🇺 Russia → ロシア語
	"\U0001F1F8\U0001F1E6": "AR",  // 🇸🇦 Saudi Arabia → アラビア語
	"\U0001F1F8\U0001F1EA": "SV",  // 🇸🇪 Sweden → スウェーデン語
	"\U0001F1F8\U0001F1EE": "SL",  // 🇸🇮 Slovenia → スロベニア語
	"\U0001F1F8\U0001F1F0": "SK",  // 🇸🇰 Slovakia → スロバキア語
	"\U0001F1F8\U0001F1F3": "WO",  // 🇸🇳 Senegal → ウォロフ語
	"\U0001F1F9\U0001F1ED": "TH",  // 🇹🇭 Thailand → タイ語
	"\U0001F1F9\U0001F1EF": "TG",  // 🇹🇯 Tajikistan → タジク語
	"\U0001F1F9\U0001F1F2": "TK",  // 🇹🇲 Turkmenistan → トルクメン語
	"\U0001F1F9\U0001F1F7": "TR",  // 🇹🇷 Turkey → トルコ語
	"\U0001F1F9\U0001F1FC": "ZH",  // 🇹🇼 Taiwan → 中国語
	"\U0001F1F9\U0001F1FF": "SW",  // 🇹🇿 Tanzania → スワヒリ語
	"\U0001F1FA\U0001F1E6": "UK",  // 🇺🇦 Ukraine → ウクライナ語
	"\U0001F1FA\U0001F1F8": "EN",  // 🇺🇸 USA → 英語
	"\U0001F1FA\U0001F1FF": "UZ",  // 🇺🇿 Uzbekistan → ウズベク語
	"\U0001F1FB\U0001F1E6": "LA",  // 🇻🇦 Vatican → ラテン語
	"\U0001F1FB\U0001F1F3": "VI",  // 🇻🇳 Vietnam → ベトナム語
	"\U0001F1FF\U0001F1E6": "AF",  // 🇿🇦 South Africa → アフリカーンス語
}
