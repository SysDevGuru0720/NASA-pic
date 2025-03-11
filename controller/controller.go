package controller

import (
	"log"
	"regexp"
	"time"

	"github.com/SysDevGuru0720/NASA-pic/config"
	"github.com/SysDevGuru0720/NASA-pic/model"
	"github.com/SysDevGuru0720/NASA-pic/util"
	"github.com/kataras/iris/v12"
)

func ShowPic(ctx iris.Context) {
	dateParam := ctx.FormValue("date")

	var dateStr string
	if dateParam == "" {
		date := time.Now()
		dateStr = date.Format("2006-01-02")
	} else {
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			log.Printf("unipic: time converting error: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		dateStr = date.Format("2006-01-02")
	}

	pic := model.Picture{
		Date: dateStr,
	}

	err := pic.GetPic()
	if err != nil {
		log.Printf("unipic: getting picture failure: %v\n", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	rxPat := regexp.MustCompile(`^(http(s?):)([/|.|\w|\s|-])*\.(?:jpg|gif|png)$`)
	if !rxPat.MatchString(pic.URL) {
		err = util.ParseTemplate(config.Global.Config.VideoTempath, pic)
		if err != nil {
			log.Printf("unipic: generating html failure: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
	} else {
		err = util.ParseTemplate(config.Global.Config.ImageTemPath, pic)
		if err != nil {
			log.Printf("unipic: generating html failure: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.ServeFile(config.Global.Config.IndexPath)
}
