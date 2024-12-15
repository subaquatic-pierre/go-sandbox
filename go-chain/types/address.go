package types

type Address [20]uint8

func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		panic("number of bytes must be 20")
	}

	return Address(b)
}
