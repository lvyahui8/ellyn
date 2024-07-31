package collections

const slotShift = 6

type bitMap struct {
	slots    []uint64
	cap      uint
	size     int
	_padding [24]byte
}

func NewBitMap(cap uint) *bitMap {
	return &bitMap{
		slots: make([]uint64, ((uint64(cap)-1)>>slotShift)+1),
		cap:   cap,
	}
}

func (m *bitMap) Set(pos uint) {
	m.checkPos(pos)
	old := m.slots[pos>>slotShift]
	m.slots[pos>>slotShift] |= uint64(1) << pos
	if m.slots[pos>>slotShift] != old {
		m.size++
	}
}

func (m *bitMap) Get(pos uint) bool {
	m.checkPos(pos)
	return m.slots[pos>>slotShift]&(uint64(1)<<pos) != 0
}

func (m *bitMap) Clear(pos uint) {
	m.checkPos(pos)
	old := m.slots[pos>>slotShift]
	m.slots[pos>>slotShift] &= ^(uint64(1) << pos)
	if m.slots[pos>>slotShift] != old {
		m.size--
	}
}

func (m *bitMap) Size() int {
	return m.size
}

func (m *bitMap) checkPos(pos uint) {
	if pos >= m.cap {
		panic("Position exceeds bitmap length")
	}
}
