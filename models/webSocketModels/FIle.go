package webSocketModels

type FileCreate struct {
	Name string
	Hash string
	Size int64
}

type FileInfo struct {
	Name     string
	Hash     string
	Size     int64
	RealSize int64
}

type FileToken struct {
	Token string
}

type FileUpload struct {
	Token     string
	ByteShift int64
}
