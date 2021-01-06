package myRPC

const (
	EDCODE = "gob"
)

type Config int

func (conf *Config) GetEdcodeConf() string  {
	return EDCODE
}
