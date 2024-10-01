package model

import "time"

type Model interface {
	Id() uint
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type baseModel struct {
	id        uint
	createdAt time.Time
	updatedAt time.Time
}

func (b baseModel) Id() uint {
	return b.id
}

func (b baseModel) CreatedAt() time.Time {
	return b.createdAt
}

func (b baseModel) UpdatedAt() time.Time {
	return b.updatedAt
}
