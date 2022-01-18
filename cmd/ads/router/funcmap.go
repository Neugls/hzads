package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"time"

	"hz.code/hz/golib/language"
)

func getFunctionMaps() template.FuncMap {
	return template.FuncMap{
		"html": func(str string) template.HTML {
			return template.HTML(str)
		},
		"i18n":    language.I18n,
		"i18nDef": language.I18nDef,
		"i18nf":   language.I18nf,
		"jsonEncode": func(source interface{}) string {
			if bytes, err := json.Marshal(source); err == nil {
				return string(bytes)
			}
			return ""
		},
		"base64Encode": func(str string) string {
			return base64.StdEncoding.EncodeToString([]byte(str))
		},
		"htmlSafe": func(text string) template.HTML {
			return template.HTML(text)
		},
		"subStr": func(str string, lens int) string {
			if len([]rune(str)) > lens {
				return string([]rune(str)[:lens]) + "..."
			}
			return str
		},
		"add": func(a, b int) int {
			return a + b
		},

		"dec": func(a, b int) int {
			return a - b
		},
		"div": func(a, b int) int {
			return a / b
		},
		"formatTime": func(unix int64, timezone, format string) string {
			if unix == 0 {
				unix = time.Now().Unix()
			}
			if timezone == "" {
				timezone = "Asia/Chongqing"
			}
			local, err := time.LoadLocation(timezone)
			if err != nil {
				log.Printf("convert unix time:%d to format %s at timezone:%s failed: %s\n", unix, format, timezone, err.Error())
				local = time.Local
			}
			return time.Unix(unix, 0).In(local).Format(format)
		},
		"fmtTime": func(format string, t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format(format)
		},
		"isTimeZero": func(t time.Time) bool {
			return t.IsZero()
		},
		"toFixed": func(fix int, f float32) string {
			return fmt.Sprintf("%."+fmt.Sprintf("%d", fix)+"f", f)
		},
	}
}
