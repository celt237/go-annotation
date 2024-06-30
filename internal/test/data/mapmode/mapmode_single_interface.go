package mapmode

import "github.com/celt237/go-annotation/internal/test/data"

// InterfaceOne  test
// @annotation(id="1", name="test")
type InterfaceOne interface {
	// Method1  test1
	// @annotation
	Method1(a1 *data.A1, a2 *data.A2, a3 *data.A3) error

	// Method2  test2
	// @annotation(name="test")
	Method2(a2 data.A2) (a3 *data.A3, err error)

	// Method3  test3
	// @annotation(name="test")
	Method3() (data.A1, error)

	// Method4  test4
	// @annotation(name="test", des="test2")
	Method4(a3 *data.A3)
}
