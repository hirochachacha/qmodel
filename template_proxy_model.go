package qmodel

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/therecipe/qt/core"
)

type templateProxyModel struct {
	*core.QIdentityProxyModel

	Header []*template.Template
	Data   []*template.Template

	buf bytes.Buffer
}

func NewTemplateProxyModel(src core.QAbstractItemModel_ITF, columnHeader []string, columnData []string) core.QAbstractItemModel_ITF {
	tmpl := template.New("main")
	tmpl.Funcs(template.FuncMap{"col": func(item interface{}, col int) (interface{}, error) {
		switch item := item.(type) {
		case *core.QIdentityProxyModel: // header
			return toValue(item.HeaderDataDefault(col, core.Qt__Horizontal, int(core.Qt__DisplayRole))), nil
		case *core.QModelIndex: // data
			return toValue(item.Model().Data(item.Sibling(item.Row(), col), int(core.Qt__DisplayRole))), nil
		default:
			return nil, fmt.Errorf("col: unsupported type")
		}
	}})

	header := make([]*template.Template, len(columnHeader))
	for i, s := range columnHeader {
		header[i] = template.Must(tmpl.New(fmt.Sprint("header[%d]", i)).Parse(s))
	}

	data := make([]*template.Template, len(columnData))
	for i, s := range columnData {
		data[i] = template.Must(tmpl.New(fmt.Sprint("data[%d]", i)).Parse(s))
	}

	m := &templateProxyModel{
		QIdentityProxyModel: core.NewQIdentityProxyModel(nil),
		Header:              header,
		Data:                data,
	}

	m.SetSourceModel(src)

	m.ConnectColumnCount(m.columnCount)
	m.ConnectHeaderData(m.headerData)
	m.ConnectData(m.data)

	return m
}

func (m *templateProxyModel) columnCount(parent *core.QModelIndex) int {
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

		m.buf.Reset()

		m.Header[section].Execute(&m.buf, m.QIdentityProxyModel)

		return core.NewQVariant14(m.buf.String())
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
