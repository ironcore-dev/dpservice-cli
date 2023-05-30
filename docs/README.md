# Documentation - dpservice-cli

## All available commands can be found [here](/docs/commands/).

To generate current command tree add this at the start of main.go:
```
	err := doc.GenMarkdownTree(cmd.Command(), "/tmp/")
	if err != nil {
		log.Fatal(err)
	}
```
run the program once and then remove this code.
This will generate a whole series of files, one for each command in the tree, in the directory specified (in this case "/tmp/").

Cobra command Markdown [docs](https://github.com/spf13/cobra/blob/main/doc/md_docs.md)

## Development details can be found [here](/docs/development/)