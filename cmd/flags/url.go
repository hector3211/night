package flags

type UrlConnection string

func (f UrlConnection) String() string {
	return string(f)
}

func (f *UrlConnection) Type() string {
	return "url"
}

func (f *UrlConnection) Set(value string) error {
	*f = UrlConnection(value)
	return nil
}
