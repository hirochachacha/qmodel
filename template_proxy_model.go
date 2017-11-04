package qmodel

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/therecipe/qt/core"
)

type templateProxyModel struct {
	*core.QIdentityProxyModel

	Header []string
	Data   []*template.Template

	buf bytes.Buffer
}

func NewTemplateProxyModel(src core.QAbstractItemModel_ITF, columnHeader []string, columnDataFormat []string, funcs template.FuncMap) core.QAbstractItemModel_ITF {
	fs := defaultFuncs()
	for k, v := range funcs {
		fs[k] = v
	}

	tmpl := template.New("main")
	tmpl.Funcs(fs)

	data := make([]*template.Template, len(columnDataFormat))
	for i, s := range columnDataFormat {
		data[i] = template.Must(tmpl.New(fmt.Sprint("data[%d]", i)).Parse(s))
	}

	m := &templateProxyModel{
		QIdentityProxyModel: core.NewQIdentityProxyModel(nil),
		Header:              columnHeader,
		Data:                data,
	}

	m.SetSourceModel(src)

	m.ConnectColumnCount(m.columnCount)
	m.ConnectHeaderData(m.headerData)
	m.ConnectData(m.data)

	return m
}

func (m *templateProxyModel) columnCount(parent *core.QModelIndex) int {
	if len(m.Header) < len(m.Data) {
		return len(m.Data)
	}
	return len(m.Header)
}

func (m *templateProxyModel) headerData(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}
	switch orientation {
	case core.Qt__Horizontal:
		if !(0 <= section && section < len(m.Header)) {
			return core.NewQVariant()
		}
		return core.NewQVariant14(m.Header[section])
	case core.Qt__Vertical:
		return core.NewQVariant7(section)
	}
	return core.NewQVariant()
}

func (m *templateProxyModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if core.Qt__ItemDataRole(role) != core.Qt__DisplayRole {
		return core.NewQVariant()
	}

	col := index.Column()

	if !(0 <= col && col < len(m.Data)) {
		return core.NewQVariant()
	}

	m.buf.Reset()

	m.Data[col].Execute(&m.buf, m.MapToSource(index))

	return core.NewQVariant14(m.buf.String())
}

func defaultFuncs() template.FuncMap {
	var cols map[string]int

	return template.FuncMap{
		"data": func(index *core.QModelIndex, col interface{}) (interface{}, error) {
			switch col := col.(type) {
			case int:
				return toValue(index.Sibling(index.Row(), col).Data(int(core.Qt__DisplayRole))), nil
			case string:
				if cols == nil {
					m := index.Model()
					ncols := m.ColumnCount(index.Parent())
					cols = make(map[string]int, ncols)
					for i := 0; i < ncols; i++ {
						cols[m.HeaderData(i, core.Qt__Horizontal, int(core.Qt__DisplayRole)).ToString()] = i
					}
				}
				if i, ok := cols[col]; ok {
					return toValue(index.Sibling(index.Row(), i).Data(int(core.Qt__DisplayRole))), nil
				}
			}
			return nil, fmt.Errorf("data: unsupported type")
		},
		"row": func(index *core.QModelIndex) int {
			return index.Row()
		},
		"col": func(index *core.QModelIndex) int {
			return index.Column()
		},
	}
}
