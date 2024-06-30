package arraymode

import "github.com/celt237/go-annotation/internal/test/data"

// InterfaceTwo  test
// @annotation
type InterfaceTwo interface {
	// Method1  test1
	// @annotation test
	Method1(a1 *data.A1, a2 *data.A2, a3 *data.A3) error

	// Method2  test2
	// @annotation test
	Method2(a2 data.A2) (a3 *data.A3, err error)

	// Method3  test3
	// @annotation test
	Method3() (data.A1, error)

	// Method4  test4
	// @annotation test  test2
	Method4(a3 *data.A3)
}

// StructTwo  test
// @annotation
type StructTwo struct {
}

// Method1  test1
// @annotation test
func (s *StructTwo) Method1(a1 *data.A1, a2 *data.A2, a3 *data.A3) error {
	return nil
}

// Method2  test2
// @annotation test
func (s *StructTwo) Method2(a2 data.A2) (a3 *data.A3, err error) {
	return nil, nil
}

// Method3  test3
// @annotation test
func (s *StructTwo) Method3() (data.A1, error) {
	return data.A1{}, nil
}

// Method4  test4
// @annotation test  test2
func (s *StructTwo) Method4(a3 *data.A3) {
	return
}
