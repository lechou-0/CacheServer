/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-16 21:08:23
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-28 00:16:56
 * @Description:
 */

package cacheserver

import (
	"log"
	"strconv"
	"sync"

	pb "wsh/server/proto"
	"wsh/server/util"

	uuid "github.com/nu7hatch/gouuid"
	zmq "github.com/pebbe/zmq4"
)

/**
 * @Author: lechoux lechoux@qq.com
 * @description: clear push socketcache. Maybe it can be optimized with LRU
 * @return {*}
 */

type CacheServer struct {
	Id             string
	ResponsePort   int
	StorageManager string
	ReadCache      map[string]*pb.KeyValuePair
	ReadCacheLock  *sync.RWMutex
	pushers        *util.SocketCache
}

func (s *CacheServer) cache_get_bind_proc() string {
	return "ipc://requests/get"
}

func (s *CacheServer) cache_set_bind_proc() string {
	return "ipc://requests/put"
}

func (s *CacheServer) cache_bind_node() string {
	return "tcp://*:" + strconv.Itoa(s.ResponsePort)
}

/**
 * @Author: lechoux lechoux@qq.com
 * @description: traverse socket and process W/R request
 * @param {[]zmq.Polled} sockets
 * @param {zmq.Socket} get_puller
 * @param {zmq.Socket} set_puller
 * @param {zmq.Socket} pusher
 * @return {*}
 */
func (s *CacheServer) execWR(sockets []zmq.Polled, get_puller *zmq.Socket, set_puller *zmq.Socket) {
	for _, polled := range sockets {
		switch socket := polled.Socket; socket {
		case get_puller:
			getRequest, err := util.RecvMsg(socket)
			if err != nil {
				log.Printf("Unexpected error while receiving get message: %v", err)
				break
			}
			getResponse := &pb.KeyResponse{}
			getResponse.Key = getRequest.GetKey()
			// get value from cache
			s.ReadCacheLock.RLock()
			val, ok := s.ReadCache[getRequest.GetKey()]
			s.ReadCacheLock.RUnlock()

			if ok {
				getResponse.Value = val.GetValue()
				getResponse.Error = pb.CacheError_NO_ERROR
			} else {
				getResponse.Error = pb.CacheError_KEY_DNE
			}

			// send msg by pushercache or new pusher
			endpoint := getRequest.GetResponseAddress()
			err = s.pushers.SendMsg(getResponse, zmq.PUSH, endpoint)
			if err != nil {
				log.Printf("Unexpected error while sending get message: %v", err)
			}

		case set_puller:
			setRequest, err := util.RecvMsg(socket)
			if err != nil {
				log.Printf("Unexpected error while receiving set message: %v", err)
				break
			}

			kv := &pb.KeyValuePair{Value: setRequest.GetValue()}
			s.ReadCacheLock.RLock()
			s.ReadCache[setRequest.GetKey()] = kv
			s.ReadCacheLock.RUnlock()

			setResponse := &pb.KeyResponse{}
			setResponse.Error = pb.CacheError_NO_ERROR

			// send msg by pushercache or new pusher
			endpoint := setRequest.GetResponseAddress()
			err = s.pushers.SendMsg(setResponse, zmq.PUSH, endpoint)
			if err != nil {
				log.Printf("Unexpected error while sending get message: %v", err)
			}
		}
	}
}

func (s *CacheServer) Run() {

	zctx, err := zmq.NewContext()
	if err != nil {
		log.Fatalf("Unexpected error while creating zmq context: %v", err)
	}
	s.pushers.Ctx = zctx
	defer zctx.Term()

	get_puller, _ := zctx.NewSocket(zmq.PULL)
	defer get_puller.Close()
	get_puller.Bind(s.cache_get_bind_proc())

	set_puller, _ := zctx.NewSocket(zmq.PULL)
	defer set_puller.Close()
	set_puller.Bind(s.cache_set_bind_proc())

	socketPoller := zmq.NewPoller()
	socketPoller.Add(get_puller, zmq.POLLIN)
	socketPoller.Add(set_puller, zmq.POLLIN)

	for {
		sockets, err := socketPoller.Poll(-1)
		if err != nil {
			log.Printf("Unexpected error while processing poller: %v", err)
		}
		go s.execWR(sockets, get_puller, set_puller)
	}
}

func NewServer() *CacheServer {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Unexpected error while generating UUID: %v", err)
	}
	server := &CacheServer{
		Id:            uid.String(),
		ResponsePort:  5555,
		ReadCache:     map[string]*pb.KeyValuePair{},
		ReadCacheLock: &sync.RWMutex{},
		pushers:       &util.SocketCache{},
	}
	log.Println("init success")
	return server
}
