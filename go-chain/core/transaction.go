package core

import "io"

type Transaction struct{}

func (tx *Transaction) EncodeBinary(r io.Writer) error {
	return nil
}
func (tx *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
