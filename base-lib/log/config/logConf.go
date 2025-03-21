package config

// Config file entity struct
type LogConf struct {
	Encoding string
	Level    string
	//Development      bool
	//OutputPaths      []string `yaml:"outputPaths"`
	//ErrorOutputPaths []string `yaml:"errorOutputPaths"`
	//EncoderConfig    EncoderConfigStruct `yaml:"encoderConfig"`
	//InitialFields    InitialFieldStruct  `yaml:"initialFields"`
	Rotate       RotateStruct
	FilterTopic  string        `yaml:"filterTopic"`
	AccessRotate *RotateStruct `yaml:"accessRotate"`
}

type EncoderConfigStruct struct {
	MessageKey    string `yaml:"messageKey"`
	LevelKey      string `yaml:"levelKey"`
	TimeKey       string `yaml:"timeKey"`
	CallerKey     string `yaml:"callerKey"`
	StacktraceKey string `yaml:"stacktraceKey"`
	// LineEnding      string `yaml:"lineEnding"`
	// LevelEncoder    string `yaml:"levelEncoder"`
	// TimeEncoder     string `yaml:"timeEncoder"`
	// DurationEncoder string `yaml:"durationEncoder"`
	// CallerEncoder   string `yaml:"callerEncoder"`
}

type InitialFieldStruct struct {
	App string
	Env string
}

type RotateStruct struct {
	Filename   string
	Maxsize    int
	Maxage     int
	Maxbackups int
	Localtime  bool
	Compress   bool
}
