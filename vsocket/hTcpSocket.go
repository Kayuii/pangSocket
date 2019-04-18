package vsocket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/Kayuii/pangSocket/util"
)

// CONSTHEADER : 包头 CONSTHEADERLENGTH : 包长 CONSTMLENGTH : 数据长
const (
	CONSTHEADER       = "Header"
	CONSTHEADERLENGTH = 6
	CONSTMLENGTH      = 4
)

type TCPSocket struct {
}

func newTCPSocket() *TCPSocket {
	tsocket := &TCPSocket{}
	return tsocket
}

// tcp接收消息
func (ts *TCPSocket) ConnHandle(msf *serv, sess *session) {
	defer func() {
		msf.SessionMaster.DelSessionByID(sess.ID)
		//调用断开链接事件
		msf.EventPool.OnClose(sess.ID)
	}()
	var errs error
	tempBuff := make([]byte, 0)
	readBuff := make([]byte, 14)
	data := make([]byte, 20)
	for {
		n, err := sess.Con.Read(readBuff)
		//更新接收时间
		sess.updateTime()
		if err != nil {
			return
		}
		tempBuff = append(tempBuff, readBuff[:n]...)
		tempBuff, data, errs = ts.Depack(tempBuff)
		if errs != nil {
			log.Println(errs)
			return
		}

		if len(data) == 0 {
			continue
		}
		//把请求的到数据转化为map
		requestData := util.String2Map(string(data))
		if msf.Hook(sess.ID, requestData) == false {
			return
		}
	}
}

// 封包
func (ts *TCPSocket) Pack(message []byte) []byte {
	return append(append([]byte(CONSTHEADER), ts.IntToBytes(len(message))...), message...)
}

// 解包
func (ts *TCPSocket) Depack(buff []byte) ([]byte, []byte, error) {
	length := len(buff)
	//如果包长小于header 就直接返回 因为接收的数据不完整
	if length < CONSTHEADERLENGTH+CONSTMLENGTH {
		return buff, nil, nil
	}

	//如果header不是 指定的header 说明此数据已经被污染 直接返回错误
	if string(buff[:CONSTHEADERLENGTH]) != CONSTHEADER {
		return []byte{}, nil, errors.New("header is not safe")
	}

	msgLength := ts.BytesToInt(buff[CONSTHEADERLENGTH : CONSTHEADERLENGTH+CONSTMLENGTH])
	if length < CONSTHEADERLENGTH+CONSTMLENGTH+msgLength {
		return buff, nil, nil
	}

	data := buff[CONSTHEADERLENGTH+CONSTMLENGTH : CONSTHEADERLENGTH+CONSTMLENGTH+msgLength]
	buffs := buff[CONSTHEADERLENGTH+CONSTMLENGTH+msgLength:]
	return buffs, data, nil
}

func (ts *TCPSocket) IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func (ts *TCPSocket) BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
