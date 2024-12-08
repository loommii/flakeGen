package flakegen

import (
	"errors"
	"sync"
	"time"
)

var (
	defaultStartUnixMilli int64 = 1288834974657 // 默认的启动时间为 Twitter 雪花算法的起始时间（2010年11月4日 01:42:54 UTC），单位为毫秒
	defaultNodeBits       uint8 = 5             // 默认的机器ID偏移
	defaultserviceBits    uint8 = 5             // 默认的业务ID偏移
	defaultNumIDBits      uint8 = 12            // 默认的序列ID偏移
)

type Node struct {
	mu           sync.Mutex
	nowUnixMilli int64 // 记录上一次的时间戳

	startUnixMilli int64 // 启动时间

	numID    uint16 // 序列ID
	maxNumID uint16

	nodeID    uint8 // 机器ID
	maxNodeID uint8

	serviceID    uint8 // 业务ID
	maxServiceID uint8

	nodeBits    uint8
	serviceBits uint8
	numIDBits   uint8
}

func NewNode(nodeID, serviceID uint8) (*Node, error) {
	node := &Node{
		startUnixMilli: defaultStartUnixMilli,
		nodeID:         nodeID,
		serviceID:      serviceID,

		nodeBits:    defaultNodeBits,
		serviceBits: defaultserviceBits,
		numIDBits:   defaultNumIDBits,

		maxNodeID:    1<<(defaultNodeBits) - 1,
		maxServiceID: 1<<(defaultserviceBits) - 1,
		maxNumID:     1<<(defaultNumIDBits) - 1,

		nowUnixMilli: time.Now().UnixMilli(),
	}
	if node.nodeID > node.maxNodeID {
		return nil, errors.New("机器ID越界")
	}
	if node.serviceID > node.maxServiceID {
		return nil, errors.New("业务ID越界")
	}
	if node.startUnixMilli > node.nowUnixMilli {
		return nil, errors.New("启动时间大于当前时间")
	}
	return node, nil
}
func (n *Node) GetID() (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	timeNow := time.Now().UnixMilli()

	if timeNow == n.nowUnixMilli {
		n.numID++
	} else {
		n.numID = 0
	}
	if n.nodeID > n.maxNodeID {
		return 0, errors.New("序列ID越界")
	}

	id := (timeNow-n.startUnixMilli)<<int64(n.nodeBits)<<int64(n.serviceBits)<<int64(n.numIDBits) | int64(n.nodeBits)<<int64(n.serviceBits)<<int64(n.numIDBits) | int64(n.serviceID)<<int64(n.numIDBits) | n.nowUnixMilli
	return id, nil
}
