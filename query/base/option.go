package base

import "io"

type OptionState struct {
	base IO
}

var DefaultOption OptionState = OptionState{
	base: NewBaseImpl(nil),
}

type Option func(*OptionState)

func SetDefaultBase(name string) Option {

	return func(opt *OptionState) {
		switch name {
		case "NoLayer":
			opt.base = &NoLayer{}
			return
		case "DobuleLayer":
			opt.base = &DoubleLayer{}
			return
		}
		opt.base = NewBaseImpl(nil)
	}
}

// NewBase initialize Base struct via buffer(buf)
func NewBase(bytes []byte) IO {
	return DefaultOption.base.NewFromBytes(bytes)
}

func NewBaseByIO(rio io.Reader, cap int) IO {
	bImpl := NewBaseImplByIO(rio, cap, nil)
	return DefaultOption.base.New(bImpl)
}

type GlobalConfig struct {
	useNewMovOfftToTable bool
}

var CurrentGlobalConfig GlobalConfig = GlobalConfig{useNewMovOfftToTable: true}

type OptGlobalConf func(*GlobalConfig)

func OptUseMovOff(t bool) OptGlobalConf {

	return func(g *GlobalConfig) {
		g.useNewMovOfftToTable = t
	}

}
