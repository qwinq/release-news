package main

import (
	_ "classOne/models"
	_ "classOne/routers"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.Run()
}
func HandlePrePage(data int)(string){
	pageIndex:=data-1
	if pageIndex<=1{
		pageIndex=1
	}
	return strconv.Itoa(pageIndex)
}
func HandleNextPage(data int)(string){
	pageIndex:=data+1
	return strconv.Itoa(pageIndex)
}
