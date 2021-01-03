package golib

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// 一致性hash实现

//声明新切片类型
type Uints []uint32

//返回切片成都
func (x Uints) Len() int {
	return len(x)
}

//比较两个数的大小
func (x Uints) Less(i, j int) bool {
	return x[i] > x[j]
}

//切片中连个值的交换
func (x Uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//当hash环上没有数据的时候
var ErrMessage = errors.New("hash is empty")

type Consistent struct {
	//hash环，key为hash值，值存放节点的key
	Circle map[uint32]string
	//已经排序的hash环
	SortedHashes Uints
	//虚拟节点个数，节点较少的时候平衡hash环
	VirtualNode int
	//读写锁
	sync.RWMutex
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}
func (c *Consistent) Delete(element string) {
	c.Lock()
	defer c.Unlock()
	c.delete(element)
}

func NewConsistent() *Consistent {
	return &Consistent{
		Circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

//自动生成key值
func (c *Consistent) GenerateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) HashKey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var search [64]byte
		copy(search[:], key)
		return crc32.ChecksumIEEE(search[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//更新排序方便查找
func (c *Consistent) UpdateSortedHashes() {
	hashes := c.SortedHashes[:0]
	//判断切片容量是否过大，如果过大重置
	if cap(c.SortedHashes)/(c.VirtualNode*4) > len(c.Circle) {
		hashes = nil
	}
	for k := range c.Circle {
		hashes = append(hashes, k)
	}
	//对所有hash值进行排序,方便查找
	sort.Sort(hashes)
	c.SortedHashes = hashes
}

// 添加节点
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.Circle[c.HashKey(c.GenerateKey(element, i))] = element
	}
	//更新排序
	c.UpdateSortedHashes()
}

//删除节点
func (c *Consistent) delete(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.Circle, c.HashKey(c.GenerateKey(element, i)))
	}
	//更新排序
	c.UpdateSortedHashes()
}

//顺时针查找附近的节点
func (c *Consistent) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.SortedHashes[x] > key
	}
	//使用二分查找算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.SortedHashes), f)
	// 闭合hash环
	if i > len(c.SortedHashes) {
		i = 0
	}
	return i
}

// 根据数据标示获取最近的服务器节点信息
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.Unlock()
	//如果为0返回错误
	if len(c.Circle) == 0 {
		return "", ErrMessage
	}
	//计算hash值
	key := c.HashKey(name)
	i := c.search(key)
	return c.Circle[c.SortedHashes[i]], nil
}
