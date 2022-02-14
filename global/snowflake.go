package global

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	waitTime = 1
	waitNum  = 1
)

var (
	//Epoch时间为 2022-01-10 16:14:46
	Epoch        int64 = 1641802486311 //毫秒数
	NodeBits     uint8 = 10            //机器节点位数
	SequenceBits uint8 = 12            //序列号位数

	mu sync.Mutex

	nodeMax     int64 = -1 ^ (-1 << NodeBits)     //机器数最大值
	nodeMask          = nodeMax << SequenceBits   //机器数位偏移量
	sequenceMax int64 = -1 ^ (-1 << SequenceBits) //序列号最大数量
	timeShift         = NodeBits + SequenceBits   //时间戳左移位数
	nodeShift         = SequenceBits              //机器数左移位数
)

type Node struct {
	mu       sync.Mutex
	epoch    time.Time
	ctx      context.Context
	time     int64
	node     int64
	sequence int64

	nodeMax     int64
	nodeMask    int64
	sequenceMax int64
	timeShift   uint8
	nodeShift   uint8
}

type ID int64

func NewNode(node int64, ctx context.Context) (*Node, error) {
	mu.Lock()
	defer mu.Unlock()
	nodeMax = -1 ^ (-1 << NodeBits)         //重新赋值
	nodeMask = nodeMax << SequenceBits      //机器数位左移到相应位置
	sequenceMax = -1 ^ (-1 << SequenceBits) //序列号最大数量
	timeShift = NodeBits + SequenceBits     //时间戳左移位数
	nodeShift = SequenceBits                //机器数左移位数

	tmp := int64(-1 ^ (-1 << NodeBits))
	if node < 0 || node > tmp {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(tmp, 64))
	}

	n := Node{}
	n.ctx = ctx
	n.node = node
	n.nodeMax = tmp
	n.nodeMask = n.nodeMax << SequenceBits
	n.sequenceMax = -1 ^ (-1 << SequenceBits)
	n.timeShift = NodeBits + SequenceBits
	n.nodeShift = SequenceBits

	var curTime = time.Now()
	// 先求当前时间和Epoch的差值，再将差值添加到curTime的单调时间上
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1e6).Sub(curTime))
	return &n, nil
}

func (n *Node) Generate() ID {
	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / 1e6
	// 获取epoch到现在的差值的毫秒数
	//比如sequenceMax是4095，seq & 4095= seq %4096
	//n.seq==0意味着获取了该同一毫秒获取了超过4096次，需要等待下一秒
	if n.time == now {
		n.sequence = (n.sequence + 1) & n.sequenceMax
		if n.sequence == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1e6
			}
		}
	} else if n.time > now {
		//时钟回调，等待waitTime秒，尝试waitNum次数
		num := 0
		for num < waitNum && n.time > now {
			time.Sleep(time.Second * waitTime)
			now = time.Since(n.epoch).Nanoseconds() / 1e6
			num++
		}
		if num == waitNum && n.time > now {
			Logger.Errorf(n.ctx, "Generate time fail, n.time > now, n.time= %d, now= %d",
				n.time, now)
		}
	} else {
		n.sequence = 0
	}
	n.time = now

	r := ID((now)<<n.timeShift | n.node<<n.nodeShift | n.sequence)
	defer n.mu.Unlock()
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}
func ParseInt64(id int64) ID {
	return ID(id)
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err
}

func (f ID) Time() int64 {
	// 获取f的时间戳，毫秒单位
	// (int64(f) >> timeShift)是当前时间-Epoch的差值
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) Node() int64 {
	//取nodeMask位数，然后右移
	return (int64(f) & nodeMask) >> nodeShift
}

func (f ID) Sequence() int64 {
	return int64(f) & sequenceMax
}
