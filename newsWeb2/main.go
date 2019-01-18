package main

import (
	_ "newsWeb2/routers"
	"github.com/astaxie/beego"
	_"newsWeb2/models"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.Run()
}

func HandlePrePage(data int) int {
	pageIndex := data - 1
	return pageIndex
}

func HandleNextPage(data int) int {
	pageIndex := data + 1
	return pageIndex
}