/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-27 22:45:05
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-27 23:45:16
 * @Description:
 */
package util

import (
	"unsafe"
	pb "wsh/server/proto"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
)

type SocketCache struct {
	Ctx         *zmq.Context
	Pushercache map[string]*zmq.Socket
}

func (socketcache *SocketCache) clear_sockets() {
	socketcache.Pushercache = make(map[string]*zmq.Socket)
}

func (socketcache *SocketCache) SendMsg(res *pb.KeyResponse, socketType zmq.Type, endpoint string) error {
	resp_string, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	if socket, ok := socketcache.Pushercache[endpoint]; ok {
		socket.Send(Bytes2String(resp_string), 1)
	} else {
		pusher, err := socketcache.Ctx.NewSocket(socketType)
		if err != nil {
			return err
		}
		pusher.Connect(endpoint)
		socketcache.Pushercache[endpoint] = pusher
		socket.Send(Bytes2String(resp_string), 1)
	}
	return nil
}

func RecvMsg(socket *zmq.Socket) (*pb.KeyRequest, error) {
	msg, _ := socket.Recv(0)
	keyRequest := &pb.KeyRequest{}
	err := proto.Unmarshal(String2Bytes(msg), keyRequest)
	if err != nil {
		return nil, err
	}

	return keyRequest, nil
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
