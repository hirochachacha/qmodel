package qmodel

import (
	"reflect"

	"github.com/therecipe/qt/core"
)

type listModel struct {
	*core.QAbstractListModel

	Value reflect.Value
}

func ListOf(v interface{}) core.QAbstractListModel_ITF {
	val := normValue(reflect.ValueOf(v))
	if !val.IsValid() {
		return nil
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil
	}

	etyp := normType(val.Type().Elem())

	switch etyp.Kind() {
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
	case reflect.String:
	case reflect.Slice:
		if etyp.Elem().Kind() != reflect.Uint8 {
			return nil
		}
	default:
		return nil
	}

	m := &listModel{
		QAbstractListModel: core.NewQAbstractListModel(nil),
		Value:              val,
	}

	m.ConnectRowCount(m.rowCount)
	m.ConnectHeaderData(m.headerData)
	m.ConnectData(m.data)

	return m
}

func (m *listModel) rowCount(parent *core.QModelIndex) int {
	return m.Value.Len()
}

func (m *listModel) headerData(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}
	switch orientation {
	case core.Qt__Horizontal:
		if section == 0 {
			return core.NewQVariant14("Value")
		}
	case core.Qt__Vertical:
		return core.NewQVariant7(section)
	}
	return core.NewQVariant()
}

func (m *listModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}

	row := index.Row()

	if !(0 <= row && row < m.rowCount(nil)) {
		return core.NewQVariant()
	}

	return toVarint(normValue(m.Value.Index(row)))
}
