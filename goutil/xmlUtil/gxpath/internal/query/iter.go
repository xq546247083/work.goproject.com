package query

import "work.goproject.com/goutil/xmlUtil/gxpath/xpath"

type Iterator interface {
	Current() xpath.NodeNavigator
}
