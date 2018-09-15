package main

import "github.com/norunners/vue"

//go:generate vueg

type Data struct {
	Message string
}

func main() {
	vue.New(
		vue.El("#app"),
		vue.Template(mainTmpl),
		vue.Data(Data{Message: "Hello WebAssembly!"}),
	)

	select {}
}
