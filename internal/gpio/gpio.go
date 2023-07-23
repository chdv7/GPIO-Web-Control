package gpio

type Gpio interface {
	ReadGPIO(path string) ([]byte, error)
	WriteGPIO(path string, data []byte) error
}
