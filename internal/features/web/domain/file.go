package web_domain

type File struct {
	buffer []byte
	name   string
}

func NewFile(buffer []byte, name string) File {
	return File{
		buffer: buffer,
		name:   name,
	}
}

func (f File) Buffer() []byte {
	return f.buffer
}

func (f File) Name() string {
	return f.name
}
