package idl

import (
	"bufio"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"sync"
	"time"
)

type PbReader struct {
	sync.Mutex
	reader       *bufio.Reader
	LastReadTime *time.Time
}

func newPbReader(r io.Reader, t *time.Time) *PbReader {
	return &PbReader {
		reader:bufio.NewReader(r),
		LastReadTime: t,
	}
}

func byteToInt32(b []byte) (int, error) {
	if len(b) < 4 {
		return 0, errors.New("byte len < 4.")
	}

	pos := 32; ret := 0
	for i := 0; i < 4; i++ {
		pos -= 8
		ret |= int(b[i]) << uint(pos)
	}
	return ret, nil
}

func int32ToByte(x int, b []byte) {
	pos := 32
	for i := 0; i < 4; i++ {
		pos -= 8
		b[i] = byte(x >> uint(pos) & 255)
	}
}

func (t *PbReader) readInt() (n int, err error) {
	data := make([]byte, 4)
	b := data

	l := 0; size := 0
	for size < 4 {
		l, err = t.reader.Read(b)
		if err != nil {
			return
		}
		size += l
		b = b[l:]
	}

	n, err = byteToInt32(data)
	if err != nil {
		return
	}
	return n, nil
}

func (t *PbReader) ReadPb(pb proto.Message) (err error) {
	t.Lock()
	defer t.Unlock()

	var pbSize int
	pbSize, err = t.readInt()
	if err != nil {
		return
	}
	if pbSize <= 0 {
		return errors.New("pbSize is 0.")
	}

	//glog.Infof("ReadPb debug, pbSize: %d\n", pbSize)
	data := make([]byte, pbSize)
	b := data

	l := 0; size := 0

	for size < pbSize {
		l, err = t.reader.Read(b)
		if err != nil {
			return
		}
		size += l
		b = b[l:]
	}

	err = proto.Unmarshal(data, pb)
	*t.LastReadTime = time.Now()
	return err
}

func (t *PbReader) readBytes(b []byte) (n int, err error) {
	return t.reader.Read(b)
}


type PbWriter struct {
	sync.Mutex
	writer       *bufio.Writer
	LastReadTime *time.Time
}

func newPbWriter(w io.Writer, t *time.Time) *PbWriter {
	return &PbWriter {
		writer:bufio.NewWriter(w),
		LastReadTime: t,
	}
}

func (t *PbWriter) writeBytes(b []byte) (err error) {
	byteSize := len(b)
	l := 0; size := 0

	for size < byteSize {
		l, err = t.writer.Write(b)
		if err != nil {
			return
		}
		size += l
		b = b[l:]
	}
	return nil
}

func (t *PbWriter) writeInt(x int) error {
	b := make([]byte, 4)
	int32ToByte(x, b)
	return t.writeBytes(b)
}

func (t *PbWriter) WritePb(pb proto.Message) (err error) {
	t.Lock()
	defer t.Unlock()

	var data []byte
	data, err = proto.Marshal(pb)

	err = t.writeInt(len(data))
	//glog.Infof("WritePb debug, pbSize: %d\n", len(data))
	if err != nil {
		return
	}

	err = t.writeBytes(data)
	if err != nil {
		return
	}

	err = t.writer.Flush()
	if err != nil {
		return
	}
	*t.LastReadTime = time.Now()
	return nil
}

func NewPbRW(rw io.ReadWriter) (reader *PbReader, writer *PbWriter) {
	t := time.Now()
	reader = newPbReader(rw, &t)
	writer = newPbWriter(rw, &t)
	return reader, writer
}

func WriteBytes(writer io.Writer, b []byte) (err error) {
	byteSize := len(b)
	l := 0; size := 0

	for size < byteSize {
		l, err = writer.Write(b)
		if err != nil {
			return
		}
		size += l
		b = b[l:]
	}
	return nil
}