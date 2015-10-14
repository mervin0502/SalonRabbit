package engine

import (
// . "mervin.me/SalonRabbit/util/log"
)
import (
	"fmt"
	"os"
	"sync"
	//"sync/atomic"
	"errors"
	//"reflect"
	"net"
	"strconv"
	"strings"
)

var (
	LocalIPofLong uint32 = 0
)

type ByteIndex [28]byte

func init() {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	arr := strings.Split(strings.Split(conn.LocalAddr().String(), ":")[0], ".")
	var j uint32
	i, err := strconv.Atoi(arr[0])
	j = uint32(i) << 24
	i, err = strconv.Atoi(arr[1])
	j += uint32(i) << 16
	i, err = strconv.Atoi(arr[2])
	j += uint32(i) << 8
	i, err = strconv.Atoi(arr[3])
	j += uint32(i)
	LocalIPofLong = j
}

//*ObjectIndex
type ObjectIndex struct {
	machine uint32
	pid     uint16
	id      uint32
	Rand    uint32
	mutex   sync.Mutex
}

//NewObjectIndex
func NewObjectIndex(h uint32) *ObjectIndex {

	return &ObjectIndex{
		machine: LocalIPofLong,
		pid:     uint16(os.Getpid()),
		id:      h,
		Rand:    0,
	}
}

//GetMachine returns
func (o *ObjectIndex) GetMachine() uint32 {
	return o.machine
}

//GetPid returns
func (o *ObjectIndex) GetPid() uint16 {
	return o.pid
}

//GetId returns
func (o *ObjectIndex) GetId() uint32 {
	return o.id
}

//GetCounter
func (o *ObjectIndex) GetRand() uint32 {
	return o.Rand
}

//SetCounter
func (o *ObjectIndex) SetRand(i uint32) {
	o.Rand = i
}

//Equals
func (o *ObjectIndex) Equals(i *ObjectIndex) bool {
	if o == i {
		return true
	}
	if i.machine == o.machine &&
		i.pid == o.pid &&
		i.id == o.id &&
		i.Rand == o.Rand {
		return true
	}
	return false
}

//ToHex return the 16-bit character of objectIndex
func (o *ObjectIndex) ToHex() string {
	var str string
	str += fmt.Sprintf("%08x", o.machine)
	str += fmt.Sprintf("%04x", o.pid)
	str += fmt.Sprintf("%08x", o.id)
	str += fmt.Sprintf("%08x", o.Rand)
	return str
}
func (o *ObjectIndex) ToByte() ByteIndex {
	var bytes ByteIndex
	var str string
	str = fmt.Sprintf("%08x", o.machine)
	copy(bytes[0:8], str)
	str = fmt.Sprintf("%04x", o.pid)
	copy(bytes[8:12], str)
	str = fmt.Sprintf("%08x", o.id)
	copy(bytes[12:20], str)
	str = fmt.Sprintf("%08x", o.Rand)
	copy(bytes[20:28], str)
	return bytes
}

//Hash return the hash value of *ObjectIndex
func (o *ObjectIndex) Hash() uint32 {
	var h uint32 = 17

	var str string
	str += fmt.Sprintf("%08x", o.machine)
	str += fmt.Sprintf("%04x", o.pid)
	str += fmt.Sprintf("%08x", o.id)
	str += fmt.Sprintf("%08x", o.Rand)
	for _, v := range str {
		h += h*23 + uint32(v)
	}
	return h
}

//LocalHash
func (o *ObjectIndex) LocalHash() uint32 {
	var h uint32 = 17

	var str string
	str += fmt.Sprintf("%08x", o.id)
	str += fmt.Sprintf("%08x", o.Rand)
	for _, v := range str {
		h += h*23 + uint32(v)
	}
	return h
}

//LocalIndex
func (o *ObjectIndex) LocalIndex() [16]byte {
	var str string
	str += fmt.Sprintf("%08x", o.id)
	str += fmt.Sprintf("%08x", o.Rand)
	var bytes [16]byte
	copy(bytes[0:15], str)
	return bytes
}

func (o *ObjectIndex) String() string {
	str := fmt.Sprintf("%d.%d.%d.%d", o.machine>>24, o.machine<<8>>24, o.machine<<16>>24, o.machine<<24>>24)
	str += fmt.Sprintf("-%d-%d", o.pid, o.Rand)
	return str
}

//UnmarshalObjectIndex

func (b *ByteIndex) Unmarshal() (*ObjectIndex, error) {
	var o *ObjectIndex

	m := b[0:8]
	p := b[8:12]
	h := b[12:20]
	c := b[20:28]
	// Log.Print(m, p, h, c)
	mUint, err1 := strconv.ParseUint(string(m), 16, 32)
	pUint, err2 := strconv.ParseUint(string(p), 16, 16)
	hUint, err3 := strconv.ParseUint(string(h), 16, 32)
	rUint, err4 := strconv.ParseUint(string(c), 16, 32)
	if err1 != nil && err2 != nil && err3 != nil && err4 != nil {
		panic(errors.New("wrong type."))
	}
	o = &ObjectIndex{
		machine: uint32(mUint),
		pid:     uint16(pUint),
		id:      uint32(hUint),
		Rand:    uint32(rUint),
	}
	return o, nil
}
func (b ByteIndex) String() string {
	return string(b[0:28])
}
func UnmarshalObjectIndex(indexStr string) (*ObjectIndex, error) {
	var o *ObjectIndex
	m := indexStr[0:8]
	p := indexStr[8:12]
	h := indexStr[12:20]
	c := indexStr[20:28]
	mUint, err1 := strconv.ParseUint(m, 16, 32)
	// print(m, "\t", p, "\t", h)
	// println(mUint)
	pUint, err2 := strconv.ParseUint(p, 16, 16)
	hUint, err3 := strconv.ParseUint(h, 16, 32)
	rUint, err4 := strconv.ParseUint(c, 16, 32)
	if err1 != nil && err2 != nil && err3 != nil && err4 != nil {
		panic(errors.New("wrong type."))
	}
	o = &ObjectIndex{
		machine: uint32(mUint),
		pid:     uint16(pUint),
		id:      uint32(hUint),
		Rand:    uint32(rUint),
	}
	return o, nil
}
