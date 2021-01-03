package icmp_tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

var icmp ICMP

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func DoPing(ip string) error {
	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	//先在buffer中写入icmp数据报求去校验和
	_ = binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	_ = binary.Write(&buffer, binary.BigEndian, icmp)

	Time, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("ip4:icmp", ip, Time)
	if err != nil {
		return fmt.Errorf("net.DialTimeout %s", err)
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("conn.Write %s", err)
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		return fmt.Errorf("conn.SetReadDeadline %s", err)
	}
	num, err := conn.Read(recvBuf)
	if err != nil {
		return fmt.Errorf("conn.Read %s", err)
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return fmt.Errorf("conn.SetReadDeadline %s", err)
	}

	if string(recvBuf[0:num]) != "" {
		return nil
	}
	return fmt.Errorf("conn.recvBuf Len Err \n")

}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16

	return uint16(^sum)
}
