package util

var (
	ADD interface{} = 1
	DEL interface{} = 2
	GET interface{} = 3
)

type channleMap struct {
	//['type','id','value',channle]
	Msq     chan *[3]interface{}
	data    map[interface{}]interface{}
	channle chan interface{}
}

// NewSafeMap : new *channleMap for map
func NewSafeMap() *channleMap {
	tem := &channleMap{}
	tem.init()
	return tem
}

func (ch *channleMap) init() {
	ch.Msq = make(chan *[3]interface{}, 10)
	ch.data = make(map[interface{}]interface{})
	ch.channle = make(chan interface{}, 0)
	go ch.run()
}

func (ch *channleMap) run() {
	for {
		select {
		case msg := <-ch.Msq:
			switch msg[0] {
			case ADD:
				ch.dataAdd(msg[1], msg[2])
			case DEL:
				ch.dataDel(msg[1])
			case GET:
				ch.dataGet(msg[1])
			}
		}
	}
}

func (ch *channleMap) msqChan(typ, id, val interface{}) *[3]interface{} {
	return &[...]interface{}{typ, id, val}
}

// 保存 或者更新元素
func (ch *channleMap) dataAdd(id, value interface{}) {
	ch.data[id] = value
}

// 删除元素
func (ch *channleMap) dataDel(id interface{}) {
	delete(ch.data, id)
}

// 获得元素
func (ch *channleMap) dataGet(id interface{}) {
	if val, exit := ch.data[id]; exit {
		ch.channle <- val
		return
	}
	ch.channle <- nil
}

//----------------------------------------------------对外接口--------------------------------
func (ch *channleMap) Add(id, value interface{}) {
	ch.Msq <- ch.msqChan(ADD, id, value)
}

func (ch *channleMap) Del(id interface{}) {
	ch.Msq <- ch.msqChan(DEL, id, nil)
}

func (ch *channleMap) Get(id interface{}) interface{} {
	ch.Msq <- ch.msqChan(GET, id, nil)
	res := <-ch.channle
	return res
}

// 获得 长度
func (ch *channleMap) GetLength() uint32 {
	return uint32(len(ch.data))
}
