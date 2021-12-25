# vueg
Command `vueg` is the Go generator for [Vue](https://github.com/norunners/vue) templates.

# Install
```bash
go get github.com/norunners/vueg
```

# Hello World!
The `mainTmpl.vue` Vue template file is read to generate Go source.
```vue
<template>
    <div>{{ Message }}</div>
</template>
```

The `main.go` file declares the `//go:generate vueg` directive.
Note, the constant `mainTmpl` is expected to be generated.
```go
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
```

Run `vueg` with the following command.
```bash
go generate
```

Finally, the `mainTmpl.go` file is generated with the `mainTmpl` constant.
```go
// The vueg command generated this source from file: mainTmpl.vue, do not edit.

package main

const mainTmpl = "<div>{{ Message }}</div>"
```

# File Watching
Below is the `vueg` file watcher setup in GoLand.
The file watcher runs `go generate` after changes are made to Go files.
Note, this configuration does not listen to changes on the Vue template files themselves.
Other IDEs that support file watchers may be configured similarly.
![file-watcher](https://user-images.githubusercontent.com/25853983/45580820-53ee0a80-b84a-11e8-9494-d9427b609f44.png)

# Styling
The `vueg` command does not currently support the `<style>` element in Vue templates.
The goal is to support scoped styling in the future.

# License
* [MIT License](LICENSE)
