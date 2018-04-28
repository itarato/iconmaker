package main

import (
	"log"

	"github.com/itarato/iconmaker/generator"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("message", "Hello world!")
		ctx.View("main.html")
	})

	app.Post("/submit", func(ctx iris.Context) {
		title := ctx.FormValue("title")
		url := ctx.FormValue("url")
		icon_file, _, err := ctx.FormFile("icon")
		if err != nil {
			log.Fatal(err)
		}

		ctx.ViewData("title", title)
		ctx.ViewData("url", url)

		package_name, err := generator.Generate(title, url, icon_file)
		if err != nil {
			ctx.View("error.html")
			return
		}

		ctx.ViewData("package_name", package_name)

		ctx.View("deliver.html")
	})

	app.Get("/download/{package:string}", func(ctx iris.Context) {
		package_name := ctx.Params().Get("package")
		ctx.SendFile("/tmp/"+package_name+".zip", package_name+".zip")
	})

	app.Run(iris.Addr(":8080"))
}
