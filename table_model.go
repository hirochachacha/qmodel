package qmodel

import (
	"reflect"

	"github.com/therecipe/qt/core"
)

type tableModel struct {
	*core.QAbstractTableModel

	Value reflect.Value
	Etype reflect.Type
}

func TableOf(v interface{}) core.QAbstractTableModel_ITF {
	val := normValue(reflect.ValueOf(v))
	if !val.IsValid() {
		return nil
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil
	}

	etyp := normType(val.Type().Elem())

	if etyp.Kind() != reflect.Struct {
		return nil
	}

	m := &tableModel{
		QAbstractTableModel: core.NewQAbstractTableModel(nil),
		Value:               val,
		Etype:               etyp,
	}

	m.ConnectColumnCount(m.columnCount)
	m.ConnectRowCount(m.rowCount)
	m.ConnectHeaderData(m.headerData)
	m.ConnectData(m.data)

	return m
}

func (m *tableModel) columnCount(parent *core.QModelIndex) int {
	return numField(m.Etype)
}

func (m *tableModel) rowCount(parent *core.QModelIndex) int {
	return m.Value.Len()
}

func (m *tableModel) headerData(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}
	switch orientation {
	case core.Qt__Horizontal:
		name, _, index := field(reflect.Zero(m.Etype), section)
		if index == -1 {
			return core.NewQVariant14(name)
		}
	case core.Qt__Vertical:
		return core.NewQVariant7(section)
	}
	return core.NewQVariant()
}

func (m *tableModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}

	row := index.Row()
	col := index.Column()

	if !(0 <= row && row < m.rowCount(nil)) {
		return core.NewQVariant()
	}

	e := normValue(m.Value.Index(row))

	_, sv, i := field(e, col)
	if i == -1 {
		return toVarint(normValue(sv))
	}

	return core.NewQVariant()
}
