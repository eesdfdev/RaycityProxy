package Stream

import (
	"RCProxy/Logger"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

type InStream struct {
	reader io.Reader
}

func NewInStream(payload []byte) *InStream {
	return &InStream{reader: bytes.NewReader(payload)}
}
func (inStream InStream) ReadByte() byte {
	var b byte
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readByte : %v", err)
		return 0xFF
	}
	return b
}
func (inStream InStream) ReadShort() int16 {
	var b int16
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readShort : %v", err)
		return -1
	}
	return b
}
func (inStream InStream) ReadUShort() uint16 {
	var b uint16
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readUShort : %v", err)
		return 0xFFFF
	}
	return b
}
func (inStream InStream) ReadInt() int32 {
	var b int32
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readInt : %v", err)
		return -1
	}
	return b
}
func (inStream InStream) ReadFloat() float32 {
	var b float32
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readFloat : %v", err)
		return -1
	}
	return b
}
func (inStream InStream) ReadUInt() uint32 {
	var b uint32
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readUInt : %v", err)
		return 0xFFFFFFFF
	}
	return b
}
func (inStream InStream) ReadInt64() int64 {
	var b int64
	err := binary.Read(inStream.reader, binary.LittleEndian, &b)
	if err != nil {
		Logger.Errorf("Error in InStream::readInt64 : %v", err)
		return 0xFFFFFFFF
	}
	return b
}
func (inStream InStream) ReadBytes(length int) []byte {
	b := make([]byte, length)
	_, err := io.ReadFull(inStream.reader, b)
	if err != nil {
		Logger.Errorf("Error in InStream::readBytes : %v", err)
		return []byte{}
	}
	return b
}
func (inStream InStream) ReadString() string {
	length := int(inStream.ReadInt())
	if length < 0 {
		Logger.Errorf("Error in InStream::readString : length is less than Zero!")
		return ""
	}
	buf := make([]uint16, length)
	err := binary.Read(inStream.reader, binary.LittleEndian, &buf)
	if err != nil {
		Logger.Errorf("Error in InStream::readString : %v", err)
		return ""
	}
	return string(utf16.Decode(buf))
}
func (inStream InStream) ReadAllBytes() []byte {
	arr, err := io.ReadAll(inStream.reader)
	if err != nil {
		Logger.Errorf("Error in InStream::readBytes : %v", err)
		return []byte{}
	}
	return arr
}

type OutStream struct {
	buffer *bytes.Buffer
}

func NewOutStream() *OutStream {
	return &OutStream{buffer: bytes.NewBuffer(nil)}
}
func NewOutStreamPacket(rtti uint32) *OutStream {
	outStream := &OutStream{buffer: bytes.NewBuffer(nil)}
	outStream.WriteUInt(rtti)
	return outStream
}
func (outStream *OutStream) WriteByte(value byte) {
	outStream.buffer.WriteByte(value)
}
func (outStream *OutStream) WriteBool(value bool) {
	if value {
		outStream.WriteByte(1)
	} else {
		outStream.WriteByte(0)
	}
}
func (outStream *OutStream) WriteShort(value int16) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteShort : %v", err)
	}
}
func (outStream *OutStream) WriteUShort(value uint16) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteUShort : %v", err)
	}
}
func (outStream *OutStream) WriteInt(value int32) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteInt : %v", err)
	}
}
func (outStream *OutStream) WriteFloat(value float32) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteFloat : %v", err)
	}
}
func (outStream *OutStream) WriteUInt(value uint32) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteUInt : %v", err)
	}
}
func (outStream *OutStream) WriteInt64(value int64) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteInt64 : %v", err)
	}
}
func (outStream *OutStream) WriteUInt64(value uint64) {
	err := binary.Write(outStream.buffer, binary.LittleEndian, value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteUInt64 : %v", err)
	}
}
func (outStream *OutStream) WriteBytes(value []byte) {
	_, err := outStream.buffer.Write(value)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteBytes : %v", err)
	}
}
func (outStream *OutStream) WriteString(value string) {
	byteStr := utf16.Encode([]rune(value))
	outStream.WriteInt(int32(utf8.RuneCountInString(value)))
	err := binary.Write(outStream.buffer, binary.LittleEndian, byteStr)
	if err != nil {
		Logger.Errorf("Error in OutStream::WriteString : %v", err)
	}
}
func (outStream *OutStream) WriteDate(target time.Time) {
	if !target.IsZero() {
		dt, _ := time.Parse("2006-01-02", "1900-01-01")
		outStream.WriteUShort(uint16(target.Sub(dt).Hours()/24) - 1)
		outStream.WriteUShort(uint16(target.Second()/4 + target.Minute()*15 + target.Hour()*60*15))
	} else {
		outStream.WriteUInt(0x0000FFFF)
	}
}
func (outStream *OutStream) WriteTime(target time.Time) {
	if !target.IsZero() {
		dt, _ := time.Parse("2006-01-02", "1900-01-01")
		outStream.WriteUShort(uint16(target.Sub(dt).Hours()/24) - 1)
		outStream.WriteUShort(uint16(target.Second()/4 + target.Minute()*15 + target.Hour()*60*15))
	} else {
		outStream.WriteUInt(0xFFFF0000)
	}
}
func (outStream *OutStream) WriteHexString(str string) {
	str = strings.ReplaceAll(str, " ", "")

	// Decode hex string to byte slice
	bytes, err := hex.DecodeString(str)
	if err != nil {
		Logger.Errorf("Error decoding hex :", err)
		return
	}
	outStream.WriteBytes(bytes)
}
func (outStream *OutStream) ToBytes() []byte {
	return outStream.buffer.Bytes()
}
func (outStream *OutStream) ToPacketBytes() []byte {
	length := uint32(len(outStream.buffer.Bytes()))
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, length)
	return append(buf, outStream.buffer.Bytes()...)
}
