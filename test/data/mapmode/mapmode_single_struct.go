package mapmode

import (
	"github.com/celt237/go-annotation/test/data"
)

// StructOne  test
// @annotation(id="1", name="test")
type StructOne struct {
}

// Method1  test1
// @annotation
func (s *StructOne) Method1(a1 *data.A1, a2 *data.A2, a3 *data.A3) error {
	return nil
}

// Method2  test2
// @annotation(name="test")
func (s *StructOne) Method2(a2 data.A2) (a3 *data.A3, err error) {
	return nil, nil
}

// Method3  test3
// @annotation(name="test")
func (s *StructOne) Method3() (data.A1, error) {
	return data.A1{}, nil
}

// Method4  test4
// @annotation(name="test", des="test2")
func (s *StructOne) Method4(a3 *data.A3) {
	return
}
