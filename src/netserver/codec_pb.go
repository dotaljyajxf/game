package netserver

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

type TCodePb struct {
	reader  io.Reader
	writer  *bufio.Writer
	readBuf *bytes.Buffer
	encBuf  *proto.Buffer
}

func NewCodecPb(reader io.Reader, writer io.Writer, aCommonLen int) *TCodePb {
	baseBuf := make([]byte, 0, aCommonLen)
	bufWriter, ok := writer.(*bufio.Writer)
	if writer != nil && !ok {
		return nil
	}
	return &TCodePb{
		reader:  reader,
		writer:  bufWriter,
		readBuf: bytes.NewBuffer(baseBuf),
		encBuf:  proto.NewBuffer(nil),
	}
}

func (this *TCodePb) CheckSum(checkSum []byte, aDataSize uint32) bool {
	body := make([]byte, 0, aDataSize)
	copy(body, this.readBuf.Bytes())
	for i := 0; i < len(body)-1; i++ {
		body[i] = body[i] ^ body[i+1]
	}

	md5Obj := md5.New()

	body[aDataSize-1] = body[aDataSize-1] ^ byte(aDataSize&0x00ff)

	for i := 0; i < len(GlobalConfig.ArrMessCode); i++ {
		messCode := GlobalConfig.ArrMessCode[i]
		codeSum := make([]byte, len(messCode)+int(aDataSize))
		copy(codeSum[0:len(messCode)], messCode)
		copy(codeSum[len(messCode):], body)

		md5Obj.Reset()
		md5Obj.Write(codeSum)

		md5Bytes := md5.Sum(nil)

		isCurrent := true
		for j := 0; j < MD5_DIGEST_LENGTH; j++ {
			if checkSum[j] != md5Bytes[j] {
				isCurrent = false
				break
			}
		}
		if isCurrent {
			return true
		}
	}

	return false
}

func (this *TCodePb) ReadBody(aReq interface{}, aDataSize uint32) error {
	this.readBuf.Reset()

	if this.readBuf.Len() < int(aDataSize) {
		l := this.readBuf.Len()
		if l == 0 {
			l = 1
		}
		for l < int(aDataSize) {
			l *= 2
		}
		if l > GlobalConfig.MaxPackageLen {
			l = GlobalConfig.MaxPackageLen
		}

		this.readBuf.Reset()
		this.readBuf.Grow(l)
	}

	buf := this.readBuf.Bytes()
	buf = buf[0:aDataSize]

	n, err := io.ReadFull(this.reader, buf)
	if err != nil {
		return err
	}
	if n != int(aDataSize) {
		return fmt.Errorf("readBody size err")
	}

	message, ok := aReq.(proto.Message)
	if !ok {
		return fmt.Errorf("read body req err")
	}

	err = proto.Unmarshal(buf, message)
	if err != nil {
		return fmt.Errorf("readBody unmarshal err")
	}
	return nil
}

func (this *TCodePb) EnCode(aData interface{}) ([]byte, error) {
	this.encBuf.Reset()
	message, ok := aData.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("change protoMessage err")
	}

	err := this.encBuf.Marshal(message)
	if err != nil {
		return nil, err
	}
	return this.encBuf.Bytes(), nil
}

func (this *TCodePb) WriteResponse(aData interface{}) error {
	buf, err := this.EnCode(aData)
	if err != nil {
		return err
	}
	_, err = this.writer.Write(buf)
	if err != nil {
		return err
	}
	return this.writer.Flush()
}

func (this *TCodePb) JustEncode(aData interface{}) ([]byte, error) {
	message, ok := aData.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("change protoMessage err")
	}
	return proto.Marshal(message)
}
