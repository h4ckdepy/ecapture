package event_processor

import (
	"bufio"
	"bytes"
	"net/http"
)

type HTTPRequest struct {
	request    *http.Request
	parserType PARSER_TYPE
	packerType PACKET_TYPE
	isDone     bool
	isInit     bool
	reader     *bytes.Buffer
}

func (this *HTTPRequest) Body() []byte {
	return nil
}

func (this *HTTPRequest) init() {

}

func (this *HTTPRequest) Name() string {
	return "HTTPRequest"
}

func (this *HTTPRequest) PacketType() PACKET_TYPE {
	return this.packerType
}

func (this *HTTPRequest) ParserType() PARSER_TYPE {
	return this.parserType
}

func (this *HTTPRequest) Write(b []byte) (int, error) {
	// 如果未初始化
	if !this.isInit {
		this.reader = bytes.NewBuffer(b)
		buf := bufio.NewReader(this.reader)
		req, err := http.ReadRequest(buf)
		if err != nil {
			return 0, err
		}
		this.request = req
		this.isInit = true
		return len(b), nil
	}

	// 如果已初始化
	l, e := this.reader.Write(b)
	if e != nil {
		return 0, e
	}

	// TODO 检测是否接收完整个包
	if false {
		this.isDone = true
	}

	return l, nil
}

func (this *HTTPRequest) detect(payload []byte) error {
	this.init()
	rd := bytes.NewReader(payload)
	buf := bufio.NewReader(rd)
	req, err := http.ReadRequest(buf)
	if err != nil {
		return err
	}
	this.parserType = PARSER_TYPE_HTTP_REQUEST
	this.request = req
	return nil
}

func (this *HTTPRequest) IsDone() bool {
	return this.isDone
}

func (this *HTTPRequest) Reset() {
	this.isDone = false
	this.isInit = false
	this.reader.Reset()
}

func (this *HTTPRequest) Display() []byte {
	// TODO 获取 http.request的body

	return this.reader.Bytes()
	return nil
}

func init() {
	hr := &HTTPRequest{}
	hr.reader = bytes.NewBuffer(nil)
	Register(hr)
}
