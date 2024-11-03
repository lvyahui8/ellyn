package collections

import "errors"

const slotShift = 6

type BitMap struct {
	slots     []uint64
	_padding0 [48]byte
	cap       uint
	_padding1 [56]byte
	size      int
	_padding2 [56]byte
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

func (m *BitMap) GetWithoutCheck(pos uint) bool {
	return m.slots[pos>>slotShift]&(uint64(1)<<pos) != 0
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

func (m *BitMap) Merge(o *BitMap) error {
	if m.cap != o.cap {
		return errors.New("inconsistent capacity cannot be merged")
	}
	for i := range m.slots {
		m.slots[i] |= o.slots[i]
	}
	return nil
}

func (m *BitMap) Size() int {
	return m.size
}

func (m *BitMap) checkPos(pos uint) {
	if pos >= m.cap {
		panic("Position exceeds bitmap length")
	}
}
