package param

func NewBOOLParam(name string) *AssemblyParam {
	return NewDefaultParam(name, BOOL)
}

func NewSINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, SINT)
}

func NewINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, INT)
}

func NewDINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, DINT)
}

func NewLINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, LINT)
}

func NewUSINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, USINT)
}

func NewUINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, UINT)
}

func NewUDINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, UDINT)
}

func NewULINTParam(name string) *AssemblyParam {
	return NewDefaultParam(name, ULINT)
}

func NewBYTEParam(name string) *AssemblyParam {
	return NewDefaultParam(name, BYTE)
}

func NewWORDParam(name string) *AssemblyParam {
	return NewDefaultParam(name, WORD)
}

func NewDWORDParam(name string) *AssemblyParam {
	return NewDefaultParam(name, DWORD)
}

func NewLWORDParam(name string) *AssemblyParam {
	return NewDefaultParam(name, LWORD)
}

func NewREALParam(name string) *AssemblyParam {
	return NewDefaultParam(name, REAL)
}

func NewLREALParam(name string) *AssemblyParam {
	return NewDefaultParam(name, LREAL)
}
