package distxn

type TestFlags struct {
	fail, timeout, deny bool
}

type TestService interface {
	Serve() error
}

type TestMicroservice struct {
	TestFlags
}

func (TestMicroservice) Serve() error {
	return nil
}

type TestStorage struct {
	TestFlags
}

func (TestStorage) Serve() error {
	return nil
}

type TestCache struct {
	TestFlags
}

func (TestCache) Serve() error {
	return nil
}
