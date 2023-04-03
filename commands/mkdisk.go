package commands

type Mkdisk struct {
	Path      string
	Fit, Unit byte
	Size      int
}

func (tmp *Mkdisk) Execute() {
	tmp.Size = tmp.ReturnSize(tmp.Size, tmp.Unit)
}
func (tmp *Mkdisk) ReturnSize(sizeI int, unit byte) int {
	var size int
	switch {
	case unit == 'k':
		size = sizeI * 1024
	case unit == 'm' || tmp.Unit == 'o':
		size = sizeI * 1024 * 1024
	}
	return size
}
func (tmp *Mkdisk) ReturnValueWithoutMarks(value string) string {
	var tmpString string
	return tmpString
}
