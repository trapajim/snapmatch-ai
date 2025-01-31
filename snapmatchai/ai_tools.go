package snapmatchai

type AIToolPropsType int

const (
	AIToolPropsTypeString AIToolPropsType = iota
	AIToolPropsTypeInt
)

type AITools struct {
	Name        string
	Description string
	Props       []AIToolProps
}

type AIToolProps struct {
	Key         string
	Type        AIToolPropsType
	Description string
}

type AIPart interface {
	isPart()
}
type Text string

func (t Text) isPart()       {}
func (t Text) isAIResponse() {}

type Blob struct {
	Data     []byte
	MIMEType string
}

func (b Blob) isPart() {}

type AIResponse interface {
	isAIResponse()
}

type FunctionArgs struct {
	Name  string
	Value any
}
type FunctionCall struct {
	FunctionName string
	Args         []FunctionArgs
}

func (f FunctionCall) isAIResponse() {}

type FunctionCallResponse struct {
	FunctionName string
	Args         []FunctionArgs
}

func (f FunctionCallResponse) isPart() {}
