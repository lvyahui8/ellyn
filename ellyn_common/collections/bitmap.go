package collections

const slotShift = 6

type BitMap struct {
	slots    []uint64
	cap      uint
	size     int
	_padding [24]byte
}

func NewBitMap(cap uint) *BitMap {
	return &BitMap{
		slots: make([]uint64, ((uint64(cap)-1)>>slotShift)+1),
		cap:   cap,
	}
}

func (m *BitMap) Set(pos uint) {
	m.checkPos(pos)
	old := m.slots[pos>>slotShift]
	m.slots[pos>>slotShift] |= uint64(1) << pos
	if m.slots[pos>>slotShift] != old {
		m.size++
	}
}

func (m *BitMap) Get(pos uint) bool {
	m.checkPos(pos)
	return m.slots[pos>>slotShift]&(uint64(1)<<pos) != 0
}

func (m *BitMap) Clear(pos uint) {
	m.checkPos(pos)
	old := m.slots[pos>>slotShift]
	m.slots[pos>>slotShift] &= ^(uint64(1) << pos)
	if m.slots[pos>>slotShift] != old {
		m.size--
	}
}

func (m *BitMap) Size() int {
	return m.size
}

func (m *BitMap) checkPos(pos uint) {
	if pos >= m.cap {
		panic("Position exceeds bitmap length")
	}
}
