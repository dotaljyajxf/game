package netserver

import (
	"bufio"
	"encoding/binary"
	"fmt"
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

type TNetCodec struct {
	rwc             io.ReadWriteCloser
	packageReader   *io.LimitedReader
	writeBuf        *bufio.Writer
	codec           *TCodePb
	closed          bool
	headSize        int
	maxPackageLen   uint32
	readHeaderBuf   []byte
	writeHeaderBuf  []byte
	checkSumRefData []byte
}

func NewNetCodec(aRwc io.ReadWriteCloser) *TNetCodec {

	limitedReader := &io.LimitedReader{R: aRwc, N: 0}

	writeBuf := bufio.NewWriter(aRwc)

	headSize := NET_HEAD_SIZE_DEF
	commonPackageLen := GlobalConfig.CommonPackageLen
	maxPackageLen := GlobalConfig.MaxPackageLen

	codec := NewCodecPb(aRwc, writeBuf, commonPackageLen)

	return &TNetCodec{
		rwc:             aRwc,
		packageReader:   limitedReader,
		writeBuf:        writeBuf,
		codec:           codec,
		closed:          false,
		headSize:        headSize,
		maxPackageLen:   uint32(maxPackageLen),
		readHeaderBuf:   make([]byte, headSize),
		writeHeaderBuf:  make([]byte, headSize),
		checkSumRefData: make([]byte, MD5_DIGEST_LENGTH),
	}
}

func (this *TNetCodec) ReadRequest(aReq interface{}, isCheckSum bool) (uint32, error) {
	rb := this.readHeaderBuf[0:this.headSize]
	n, err := io.ReadFull(this.rwc, rb)

	if err != nil {
		return 0, err
	}
	if n != this.headSize {
		return 0, fmt.Errorf("headSize unMatch %d", n)
	}

	flag := this.readHeaderBuf[NET_HEAD_FLAG_BYTE]
	if flag&HeadFlagPing > 0 {
		return 0, nil
	}

	packageLen := binary.BigEndian.Uint32(this.readHeaderBuf)

	if packageLen > this.maxPackageLen {
		return 0, fmt.Errorf("packegLen over limit %d", packageLen)
	}

	if flag&HeadFlagCheckSum > 0 {
		n, err := io.ReadFull(this.rwc, this.checkSumRefData)
		if err != nil {
			return 0, err
		}
		if n != MD5_DIGEST_LENGTH {
			return 0, fmt.Errorf("checkData length error %d", n)
		}
	}

	err = this.codec.ReadBody(aReq, packageLen)
	if err != nil {
		return 0, err
	}

	if isCheckSum {
		if !this.codec.CheckSum(this.checkSumRefData, packageLen) {
			return 0, fmt.Errorf("checkSum error")
		}
	}
	return packageLen, nil
}

func (this *TNetCodec) WriteResponse(aData interface{}) (uint32, error) {
	pkg, err := this.codec.EnCode(aData)
	if err != nil {
		return 0, err
	}

	packageLen := uint32(len(pkg))
	binary.BigEndian.PutUint32(this.writeHeaderBuf, packageLen)

	this.writeHeaderBuf[NET_HEAD_FLAG_BYTE] = 0

	n, err := this.writeBuf.Write(this.writeHeaderBuf)
	if err != nil {
		return 0, err
	}
	if n != this.headSize {
		return 0, fmt.Errorf("send reponse head length err ")
	}

	n, err = this.writeBuf.Write(pkg)
	if err != nil {
		return 0, err
	}
	if n != int(packageLen) {
		return 0, fmt.Errorf("send reponse body length err ")
	}

	return packageLen, nil
}
func (this *TNetCodec) Encode(aData interface{}) ([]byte, error) {
	return this.codec.EnCode(aData)
}
func (this *TNetCodec) JustEncode(aData interface{}) ([]byte, error) {
	return this.codec.JustEncode(aData)
}
func (this *TNetCodec) WriteQuickResponse(aData interface{}) (uint32, error) {
	pkg, err := this.codec.EnCode(aData)
	if err != nil {
		return 0, err
	}

	packageLen := uint32(len(pkg))
	binary.BigEndian.PutUint32(this.writeHeaderBuf, packageLen)

	this.writeHeaderBuf[NET_HEAD_FLAG_BYTE] = 0

	n, err := this.rwc.Write(this.writeHeaderBuf)
	if err != nil {
		return 0, err
	}
	if n != this.headSize {
		return 0, fmt.Errorf("send reponse head length err ")
	}

	n, err = this.rwc.Write(pkg)
	if err != nil {
		return 0, err
	}
	if n != int(packageLen) {
		return 0, fmt.Errorf("send reponse body length err ")
	}

	return packageLen, nil
}
func (this *TNetCodec) SendPing() error {
	packageLen := uint32(0)
	binary.BigEndian.PutUint32(this.writeHeaderBuf, packageLen)

	this.writeHeaderBuf[NET_HEAD_FLAG_BYTE] |= HeadFlagPing

	_, err := this.rwc.Write(this.writeHeaderBuf)
	return err
}
func (this *TNetCodec) Flush() error {
	return this.writeBuf.Flush()
}
func (this *TNetCodec) Close() error {
	if this.closed {
		return nil
	}
	this.closed = true
	return this.rwc.Close()
}
