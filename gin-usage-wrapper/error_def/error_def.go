package errordef

type CliError struct {
	Content   interface{} `json:"content"`
	Message   string      `json:"message,omitempty"`   //应答备注
	MessageId string      `json:"messageId,omitempty"` //消息id
	Status    int32       `json:"status"`
	Timestamp int64       `json:"timestamp,omitempty"`
	CostTime  int64       `json:"costtime,omitempty"`
}
