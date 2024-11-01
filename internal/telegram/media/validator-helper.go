package media

type mediaType int

const (
	undefinedType mediaType = iota
	photoType
	videoType
)

type size struct {
	width    int
	height   int
	fileSize int64
}

func newSize(width, height int, bytes int64) size {
	return size{
		width:    width,
		height:   height,
		fileSize: bytes,
	}
}

type Limit struct {
	min map[mediaType]size
	max map[mediaType]size
}

func NewLimit() Limit {
	return Limit{
		min: make(map[mediaType]size),
		max: make(map[mediaType]size),
	}
}

func (l *Limit) SetMinLimit(t mediaType, s size) *Limit {
	l.min[t] = s

	return l
}

func (l *Limit) SetMaxLimit(t mediaType, s size) *Limit {
	l.max[t] = s

	return l
}
