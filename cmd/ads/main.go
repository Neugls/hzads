package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"hz.code/hz/golib/language"
	"hz.code/neugls/ads/cmd/ads/router"
	"hz.code/neugls/ads/cmd/ads/ver"
	"hz.code/neugls/ads/emd"
	"hz.code/neugls/ads/internal/config"
	"hz.code/neugls/ads/internal/database"
)

func main() {
	listen := flag.String("l", "0.0.0.0:8898", "the address to listen on")
	dataDir := flag.String("d", "", "the directory to store data in, if not set, will use the AppData directory")
	locale := flag.String("locale", "en-US", "the locale to use")

	flag.Parse()
	config.Setup(*listen, *dataDir, "hzads.db", "hzads_")

	language.LoadFromString(emd.ResLanguage)
	language.SetLocale(*locale)

	//setup db
	if err := database.Setup(); err != nil {
		panic("setup database fail: " + err.Error())
	}
	defer database.Close()

	log.Printf(os.Args[0])

	log.Printf("hz.ads version %s, build time: %s, branch: %s, commit: %s\n", ver.Version, ver.Time, ver.BRANCH, ver.COMMIT)

	router.Build().Run(config.V.Listen)
}
