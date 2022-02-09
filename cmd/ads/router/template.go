package router

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
	"hz.code/hz/golib/language"
	"hz.code/neugls/ads/internal/config"
)

//HTMLRender  html render
type HTMLRender struct {
	Tmpls   map[string]*template.Template
	FuncMap template.FuncMap

	//DefaultTitle 默认页面标题
	DefaultTitle string

	//DefaultKeywords 默认页面关键字
	DefaultKeywords string

	//DefaultDescription 默认页面描述
	DefaultDescription string

	//GoogleAnalyticsID google 统计分析的ID
	GoogleAnalyticsID string

	layoutFiles []string
	base        fs.FS
}

func newHTMLRender() HTMLRender {
	r := HTMLRender{
		Tmpls:              map[string]*template.Template{},
		DefaultTitle:       config.V.AppName,
		DefaultDescription: language.I18nDef("no fear for contract", "明棠，无惧合同"),
		DefaultKeywords:    language.I18nDef("contract,contract manager", "合同,合同管理"),
		GoogleAnalyticsID:  "",
	}
	return r
}

//Init Init
func (r *HTMLRender) Init(viewsDir fs.FS, fm template.FuncMap) error {
	r.base = viewsDir
	r.FuncMap = getFunctionMaps()
	r.FuncMap = r.mergeFuncMap(fm)

	layoutFiles, includeFiles := []string{}, []string{}
	fs.WalkDir(viewsDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("html render init walk dir %s fail, err: %s", path, err)
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			name := strings.Replace(path, string(os.PathSeparator), ".", -1)
			if strings.HasPrefix(name, "layouts.") {
				layoutFiles = append(layoutFiles, path)
			} else {
				includeFiles = append(includeFiles, path)
			}
		}
		return nil
	})

	r.layoutFiles = layoutFiles

	for _, f := range includeFiles {
		// name := strings.TrimPrefix(f, root)
		name := strings.Replace(f, string(os.PathSeparator), ".", -1)
		name = strings.TrimPrefix(name, "layouts.")

		r.Tmpls[name] = template.New(filepath.Base(f))
		files := append(layoutFiles, f)
		r.Tmpls[name] = r.Tmpls[name].Funcs(r.FuncMap)
		r.Tmpls[name] = template.Must(r.Tmpls[name].ParseFS(r.base, files...))
	}

	return nil
}

func (r *HTMLRender) mergeFuncMap(fm template.FuncMap) template.FuncMap {
	for k, v := range r.FuncMap {
		fm[k] = v
	}
	return fm
}

//Instance implements the html render interface of gin
func (r HTMLRender) Instance(name string, data interface{}) render.Render {
	// if gin.IsDebugging() {
	// 	name = strings.ReplaceAll(name, ".", string(os.PathSeparator))
	// 	name = strings.Replace(name, "/html", ".html", 1)
	// 	f := path.Join(r.base, name)
	// 	fmt.Println(f + "..")
	// 	if _, e := os.Stat(f); e == nil {
	// 		t := template.New(filepath.Base(f))
	// 		fs := []string{}
	// 		fs = append(fs, r.layoutFiles...)
	// 		fs = append(fs, f)
	// 		t.Funcs(r.FuncMap)
	// 		t = template.Must(t.ParseFiles(fs...))

	// 		return render.HTML{
	// 			Template: t,
	// 			Name:     "",
	// 			Data:     data,
	// 		}
	// 	}
	// 	panic(fmt.Sprintf("get html instance fail, name: %s, data: %s", name, data))
	// }

	if t, ok := r.Tmpls[name]; ok {
		return render.HTML{
			Template: t,
			Name:     "",
			Data:     data,
		}
	}
	panic("invalid template name:" + name)
}
