package qmodel

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/therecipe/qt/core"
)

type treeModel struct {
	*core.QAbstractItemModel

	Root *treeItem
}

type treeItem struct {
	Name  string
	Value reflect.Value

	Children []*treeItem
	Parent   *treeItem
	Row      int
}

func TreeOf(v interface{}) core.QAbstractItemModel_ITF {
	val := normValue(reflect.ValueOf(v))
	if !val.IsValid() {
		return nil
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	m := &treeModel{
		QAbstractItemModel: core.NewQAbstractItemModel(nil),
		Root: &treeItem{
			Value: val,
		},
	}

	m.ConnectColumnCount(m.columnCount)
	m.ConnectRowCount(m.rowCount)
	m.ConnectHeaderData(m.headerData)
	m.ConnectData(m.data)
	m.ConnectIndex(m.index)
	m.ConnectParent(m.parent)

	return m
}

func (m *treeModel) toItem(index *core.QModelIndex) *treeItem {
	if !index.IsValid() {
		return m.Root
	}
	return (*treeItem)(index.InternalPointer())
}

func (m *treeModel) columnCount(parent *core.QModelIndex) int {
	return 2
}

func (m *treeModel) rowCount(parent *core.QModelIndex) (count int) {
	pitem := m.toItem(parent)

	if pitem.Children != nil {
		return len(pitem.Children)
	}

	defer func() {
		pitem.Children = make([]*treeItem, count)
	}()

	pval := pitem.Value
	if !pval.IsValid() {
		return 0
	}

	switch pval.Kind() {
	case reflect.Array:
		return pval.Len()
	case reflect.Slice:
		if pval.Type().Elem().Kind() == reflect.Uint8 {
			return 0
		}
		return pval.Len()
	case reflect.Struct:
		return numField(pval.Type())
	}

	return 0
}

func (m *treeModel) headerData(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}
	switch orientation {
	case core.Qt__Horizontal:
		switch section {
		case 0:
			return core.NewQVariant14("Name")
		case 1:
			return core.NewQVariant14("Value")
		}
	case core.Qt__Vertical:
		return core.NewQVariant7(section)
	}
	return core.NewQVariant()
}

func (m *treeModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}

	item := m.toItem(index)

	switch index.Column() {
	case 0:
		if item.Name != "" {
			return core.NewQVariant14(item.Name)
		}
		return core.NewQVariant14(fmt.Sprintf("[%d]", item.Row))
	case 1:
		return toVarint(item.Value)
	}

	return core.NewQVariant()
}

func (m *treeModel) index(row, col int, parent *core.QModelIndex) *core.QModelIndex {
	pitem := m.toItem(parent)

	pval := pitem.Value
	if !pval.IsValid() {
		return core.NewQModelIndex()
	}

	if !(0 <= row && row < m.rowCount(parent)) {
		return core.NewQModelIndex()
	}

	if item := pitem.Children[row]; item != nil {
		return m.CreateIndex2(row, col, uintptr(unsafe.Pointer(item)))
	}

	var name string
	var val reflect.Value

	switch pval.Kind() {
	case reflect.Array:
		val = pval.Index(row)
	case reflect.Slice:
		if pval.Type().Elem().Kind() == reflect.Uint8 {
			return core.NewQModelIndex()
		}
		val = pval.Index(row)
	case reflect.Struct:
		if col != 0 && col != 1 {
			return core.NewQModelIndex()
		}
		n, sv, index := field(pval, row)
		if index == -1 {
			name = n
			val = sv
		}
	}

	item := &treeItem{Name: name, Value: normValue(val), Parent: pitem, Row: row}

	pitem.Children[row] = item

	return m.CreateIndex2(row, col, uintptr(unsafe.Pointer(item)))
}

func (m *treeModel) parent(index *core.QModelIndex) *core.QModelIndex {
	item := m.toItem(index)
	if item.Parent == nil {
		return core.NewQModelIndex()
	}
	return m.CreateIndex2(item.Row, 0, uintptr(unsafe.Pointer(item.Parent)))
}
