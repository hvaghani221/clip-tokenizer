package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToFirst(t *testing.T) {
	list := NewList[int, int]()
	list.AddToFront(1, 1)
	assert.EqualValues(t, []int{1}, list.GetList())
	list.AddToFront(2, 1)
	assert.EqualValues(t, []int{2, 1}, list.GetList())
	list.AddToFront(3, 1)
	assert.EqualValues(t, []int{3, 2, 1}, list.GetList())
}

func TestMoveToFront(t *testing.T) {
	list := NewList[int, int]()
	one := list.AddToFront(1, 1)
	two := list.AddToFront(2, 2)
	_ = list.AddToFront(3, 3)
	list.MoveToFront(one)
	assert.EqualValues(t, []int{1, 3, 2}, list.GetList())
	list.MoveToFront(two)
	assert.EqualValues(t, []int{2, 1, 3}, list.GetList())
}

func TestRemoveLast(t *testing.T) {
	list := NewList[int, int]()
	list.AddToFront(1, 1)
	list.AddToFront(2, 1)
	list.AddToFront(3, 1)
	assert.Equal(t, []int{3, 2, 1}, list.GetList())
	list.RemoveBack()
	assert.Equal(t, []int{3, 2}, list.GetList())
	list.RemoveBack()
	assert.Equal(t, []int{3}, list.GetList())
	list.RemoveBack()
	assert.Equal(t, []int{}, list.GetList())
}

func TestLRU1(t *testing.T) {
	obj := NewLRU[int, int](2)
	obj.Put(1, 1)
	obj.Put(2, 2)
	val, found := obj.Get(1)
	assert.Equal(t, true, found)
	assert.Equal(t, 1, val)
	obj.Put(3, 3)
	val, found = obj.Get(2)
	assert.Equal(t, false, found)
	assert.Equal(t, 0, val)
	obj.Put(4, 4)
	obj.Get(1)
	assert.Equal(t, false, found)
	assert.Equal(t, 0, val)
	val, found = obj.Get(3)
	assert.Equal(t, true, found)
	assert.Equal(t, 3, val)
	val, found = obj.Get(4)
	assert.Equal(t, true, found)
	assert.Equal(t, 4, val)
}
