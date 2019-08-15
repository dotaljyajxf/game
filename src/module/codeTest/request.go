package codeTest

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
)

/**
包头长度5：

4 bytes 包体长

1 byte
    1 bit: 是否校验数据
    2 bit: ping

如果设置了校验数据，包头后会有16个字节的检验数据

*/

const (
	NET_HEAD_SIZE_DEF = 5 //默认网络包包头长度
	//	NET_HEAD_SIZE_AMF  = 8 //Amf网络包包头长度
	NET_HEAD_FLAG_BYTE = 4 //标志byte，记录 是否ping、是否需要md5校验
	MD5_DIGEST_LENGTH  = 16
)

// 标志byte中的，bit位的含义
const (
	HeadFlagNull     byte = iota
	HeadFlagCheckSum      = 1 << uint(iota-1)
	HeadFlagPing
)

type TNetCode struct {
	rwc         io.ReadWriteCloser
	readBuffer  []byte
	writeBuffer []byte
}

func NewTNetCode(rwc io.ReadWriteCloser) *TNetCode {
	return &TNetCode{
		rwc:         rwc,
		readBuffer:  make([]byte, NET_HEAD_SIZE_DEF),
		writeBuffer: make([]byte, NET_HEAD_SIZE_DEF),
	}
}

func (this *TNetCode) ReadResponse() ([]byte, error) {
	n, err := io.ReadFull(this.rwc, this.readBuffer)
	if err != nil {
		return nil, err
	}
	if n != NET_HEAD_SIZE_DEF {
		return nil, errors.New("readHead length error")
	}
	flag := this.readBuffer[NET_HEAD_FLAG_BYTE]
	if HeadFlagPing&flag > 0 {
		return nil, nil
	}

	packageLen := binary.BigEndian.Uint32(this.readBuffer)
	if packageLen == 0 {
		return nil, nil
	}
	body := make([]byte, packageLen)

	n, err = io.ReadFull(this.rwc, body)
	if err != nil {
		return nil, err
	}
	if n != int(packageLen) {
		return nil, errors.New("readBody length error")
	}

	return body, nil
}

func (this *TNetCode) WriteRequest(req []byte, isCheck bool) error {
	packegLen := len(req)

	if isCheck {
		this.writeBuffer[NET_HEAD_FLAG_BYTE] = HeadFlagCheckSum
	} else {
		this.writeBuffer[NET_HEAD_FLAG_BYTE] = HeadFlagNull
	}

	binary.BigEndian.PutUint32(this.writeBuffer, uint32(packegLen))

	n, err := this.rwc.Write(this.writeBuffer)
	if err != nil {
		return err
	}
	if n != NET_HEAD_SIZE_DEF {
		return errors.New("write Head error")
	}

	if isCheck {
		this.writeSum(req)
	}

	n, err = this.rwc.Write(req)
	if err != nil {
		return err
	}
	if n != packegLen {
		return errors.New("write req error")
	}

	return nil
}

func (this *TNetCode) writeSum(req []byte) {
	messCode := "12332"

	data := make([]byte, len(req))
	copy(data, req)

	for i := 0; i < len(req)-1; i++ {
		data[i] = data[i] ^ data[i+1]
	}
	data[len(req)-1] = data[len(req)-1] ^ byte(len(req)&0x00ff)

	messByteArr := make([]byte, len(req)+len(messCode))
	copy(messByteArr, messCode)
	copy(messByteArr[len(messCode):], data)

	md5Obj := md5.New()
	md5Obj.Reset()
	md5Obj.Write(messByteArr)
	md5Bytes := md5Obj.Sum(nil)

	if len(md5Bytes) != MD5_DIGEST_LENGTH {
		panic("md5 length not 16")
	}

	this.rwc.Write(md5Bytes)
}

type TPbCode struct {
}

func (this *TPbCode) DeCode(buf []byte, data interface{}) error {
	msg, ok := data.(proto.Message)
	if !ok {
		panic("pb data to message error")
	}
	return proto.Unmarshal(buf, msg)
}

func (this *TPbCode) EnCode(buf interface{}) ([]byte, error) {
	msg, ok := buf.(proto.Message)
	if !ok {
		panic("pb data to message error")
	}
	return proto.Marshal(msg)
}
