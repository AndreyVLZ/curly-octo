package model

import "io"

type File struct {
	id string
	io.ReadWriteCloser
}

func (f *File) ID() string { return f.id }

func NewFile(id string, file io.ReadWriteCloser) *File {
	return &File{
		id:              id,
		ReadWriteCloser: file,
	}
}

type DecFile struct {
	id string
	io.WriteCloser
}

func (f *DecFile) ID() string { return f.id }

func NewDecFile(id string, wc io.WriteCloser) *DecFile {
	return &DecFile{
		id:          id,
		WriteCloser: wc,
	}
}

type EncFile struct {
	id string
	io.ReadCloser
}

func (f *EncFile) ID() string { return f.id }

func NewEncFile(id string, rc io.ReadCloser) *EncFile {
	return &EncFile{
		id:         id,
		ReadCloser: rc,
	}
}
