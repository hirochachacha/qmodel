package qmodel

import (
	"encoding/base64"
	"reflect"
	"strings"

	"github.com/therecipe/qt/core"
)

func normType(typ reflect.Type) reflect.Type {
	for {
		kind := typ.Kind()
		if kind != reflect.Interface && kind != reflect.Ptr {
			return typ
		}
	}
}

func normValue(val reflect.Value) reflect.Value {
	for {
		kind := val.Kind()
		if kind != reflect.Interface && kind != reflect.Ptr {
			return val
		}
		val = val.Elem()
		if !val.IsValid() {
			return val
		}
	}
}

func toVarint(val reflect.Value) *core.QVariant {
	switch val.Kind() {
	case reflect.Invalid:
		return core.NewQVariant()
	case reflect.Bool:
		return core.NewQVariant11(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return core.NewQVariant9(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return core.NewQVariant10(val.Uint())
	case reflect.Float32, reflect.Float64:
		return core.NewQVariant12(val.Float())
	case reflect.String:
		return core.NewQVariant14(val.String())
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			bs := val.Bytes()
			dst := make([]byte, base64.StdEncoding.EncodedLen(len(bs)))
			base64.StdEncoding.Encode(dst, bs)
			return core.NewQVariant14(string(dst))
		}
	}

	return core.NewQVariant()
}

func toValue(v *core.QVariant) interface{} {
	switch v.Type() {
	case core.QVariant__Invalid:
		return nil
	case core.QVariant__Bool:
		return v.ToBool()
	case core.QVariant__Int:
		return v.ToInt(true)
	case core.QVariant__UInt:
		return v.ToUInt(true)
	case core.QVariant__LongLong:
		return v.ToLongLong(true)
	case core.QVariant__ULongLong:
		return v.ToULongLong(true)
	case core.QVariant__Double:
		return v.ToDouble(true)
	case core.QVariant__Char:
		return v.ToChar()
	case core.QVariant__Map:
		return v.ToMap()
	case core.QVariant__List:
		return v.ToList()
	case core.QVariant__String:
		return v.ToString()
	case core.QVariant__StringList:
		return v.ToStringList()
	case core.QVariant__ByteArray:
		return v.ToByteArray()
	case core.QVariant__BitArray:
		return v.ToBitArray()
	case core.QVariant__Date:
		return v.ToDate()
	case core.QVariant__Time:
		return v.ToTime()
	case core.QVariant__DateTime:
		return v.ToDateTime()
	case core.QVariant__Url:
		return v.ToUrl()
	case core.QVariant__Locale:
		return v.ToLocale()
	case core.QVariant__Rect:
		return v.ToRect()
	case core.QVariant__RectF:
		return v.ToRectF()
	case core.QVariant__Size:
		return v.ToSize()
	case core.QVariant__SizeF:
		return v.ToSizeF()
	case core.QVariant__Line:
		return v.ToLine()
	case core.QVariant__LineF:
		return v.ToLineF()
	case core.QVariant__Point:
		return v.ToPoint()
	case core.QVariant__PointF:
		return v.ToPointF()
	case core.QVariant__RegExp:
		return v.ToRegExp()
	case core.QVariant__RegularExpression:
		return v.ToRegularExpression()
	case core.QVariant__Hash:
		return v.ToHash()
	case core.QVariant__EasingCurve:
		return v.ToEasingCurve()
	case core.QVariant__Uuid:
		return v.ToUuid()
	case core.QVariant__ModelIndex:
		return v.ToModelIndex()
	case core.QVariant__PersistentModelIndex:
		return v.ToPersistentModelIndex()
	case core.QVariant__Font:
		return v.ToFont()
	case core.QVariant__Pixmap:
	case core.QVariant__Brush:
	case core.QVariant__Color:
		return v.ToColor()
	case core.QVariant__Palette:
	case core.QVariant__Image:
		return v.ToImage()
	case core.QVariant__Polygon:
	case core.QVariant__Region:
	case core.QVariant__Bitmap:
	case core.QVariant__Cursor:
	case core.QVariant__KeySequence:
	case core.QVariant__Pen:
	case core.QVariant__TextLength:
	case core.QVariant__TextFormat:
	case core.QVariant__Matrix:
	case core.QVariant__Transform:
	case core.QVariant__Matrix4x4:
	case core.QVariant__Vector2D:
	case core.QVariant__Vector3D:
	case core.QVariant__Vector4D:
	case core.QVariant__Quaternion:
	case core.QVariant__PolygonF:
	case core.QVariant__Icon:
		return v.ToIcon()
	case core.QVariant__SizePolicy:
	case core.QVariant__UserType:
	case core.QVariant__LastType:
	}
	return nil
}

func numField(typ reflect.Type) int {
	var count int

	n := typ.NumField()
	for i := 0; i < n; i++ {
		sf := typ.Field(i)

		tag := sf.Tag.Get("qmodel")
		if tag == "-" {
			continue
		}

		if sf.Anonymous { // embeded
			if tag == "" {
				st := sf.Type
				if st.Kind() == reflect.Ptr {
					st = st.Elem()
				}
				if st.Kind() == reflect.Struct {
					count += numField(st)
					continue
				}
			}
		}

		if sf.PkgPath == "" { // exported
			count++
		}
	}

	return count
}

func field(val reflect.Value, index int) (string, reflect.Value, int) {
	n := val.NumField()
	for i := 0; i < n; i++ {
		sf := val.Type().Field(i)

		tag := sf.Tag.Get("qmodel")
		if tag == "-" {
			continue
		}

		sv := val.Field(i)

		if sf.Anonymous { // embeded
			if tag == "" {
				if sv.Kind() == reflect.Ptr {
					sv = sv.Elem()
				}
				if sv.Kind() == reflect.Struct {
					var n1 string
					var sv1 reflect.Value
					n1, sv1, index = field(sv, index)
					if index == -1 {
						return n1, sv1, -1
					}
					continue
				}
			}
		}

		if sf.PkgPath == "" { // exported
			index--
			if index == -1 {
				name := sf.Name
				if tag := sf.Tag.Get("qmodel"); tag != "" {
					opts := strings.Split(tag, ",")
					if opts[0] != "" {
						name = opts[0]
					}
				}
				return name, sv, -1
			}
		}
	}
	return "", reflect.Value{}, index
}
