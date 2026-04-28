package launcher

type itemType int

const (
	appItem itemType = iota
	commandItem
	fileItem
)

type item struct {
	title string
	cmd   string
	typ   itemType
}

type usage struct {
	Count int
	Last  int64
}

type mode string

const (
	appMode  mode = "app"
	cmdMode  mode = "cmd"
	fileMode mode = "file"
)
