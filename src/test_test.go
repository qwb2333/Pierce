package test

import (
	"testing"
	"io/ioutil"
	"os"
	"idl"
	"github.com/golang/protobuf/proto"
	"bytes"
	"config"
)

func TestConfig(t *testing.T) {
	str := `# this is tips.
			outerIp=127.0.0.1
			outerPort=4910`
	fileName := "/tmp/testConfig.properties"
	ioutil.WriteFile(fileName, []byte(str), os.ModePerm)
	defer os.Remove(fileName)

	c := config.NewConfig(fileName)
	host, _ := c.ReadString("outerIp")
	port, _ := c.ReadInt("outerPort")
	if port != 4910 || host != "127.0.0.1" {
		t.Fatal("host or post is null.")
	}
}

func TestProtobuf(t *testing.T) {
	info := &idl.RemoteInfo {
		Ip: proto.String("127.0.0.1"),
		Port: proto.Int32(80),
	}

	data, err := proto.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}

	info2 := &idl.RemoteInfo{}
	proto.Unmarshal(data, info2)
	if *info2.Ip != "127.0.0.1" || *info2.Port != 80 {
		t.Fatal("failed.")
	}
}

func TestPbRW(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	rb, wb := idl.NewPbRW(buffer)

	info := &idl.RemoteInfo {
		Ip: proto.String("127.0.0.1"),
	}

	for i := 0; i < 1000; i++ {
		info.Port = proto.Int(i)
		err := wb.WritePb(info)
		if err != nil {
			t.Fatal(err)
		}
	}


	for i := 0; i < 1000; i++ {
		info2 := idl.RemoteInfo{}
		err := rb.ReadPb(&info2)
		if err != nil {
			t.Fatal(err)
		}

		if *info2.Ip != "127.0.0.1" || int(*info2.Port) != i {
			t.Fatal("TestPbRW failed.")
		}
	}
}

/*
如果要实现多态的话，需要interface和struct,其中struct继承了那个interface
然后对这个struct取&，其内容对应了interface，注意这时的interface并不是指针

对于既有变量，又有interface的struct，应该拆成两部分，一部分是interface，一部分是struct
如果要实现interface操作资源，那么就必须要写GetX和SetX之类的东西
如果是传递interface，是不会复制对应的struct里的内存的
如果是传递struct，则会复制内存
*/