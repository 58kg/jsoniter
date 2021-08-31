package jsoniter

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
)

// 结点处理器
type Handler struct {
	Fields  []string
	Handler func(obj interface{}, fields []string) error // 参数obj的类型可能是map[string]interface{}或[]interface{}且为fields[len(fields)-1]的父节点
}

func Traverse(obj interface{}, handlers []Handler) (err error) {
	// 确认obj可以序列化为json串
	if _, err := json.Marshal(obj); err != nil {
		return fmt.Errorf("obj can not marshal to a json string, the error in json.Marshal is %v", err)
	}
	if len(handlers) == 0 {
		return errors.New("handler should not be empty")
	}
	// 检查handlers的合法性
	for _, h := range handlers {
		if h.Handler == nil {
			return errors.New("some handler is empty")
		}
		if len(h.Fields) == 0 {
			return errors.New("some handler fields is empty")
		}
		for _, v := range h.Fields {
			if len(v) == 0 {
				return errors.New("some handler fields is empty")
			}
		}
	}

	// 创建trie树
	root := &node{children: map[string]*node{}}
	for _, h := range handlers {
		if err := root.insert(h.Fields, h.Handler); err != nil {
			return err
		}
	}

	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("fatal error: panic in dfs, [Err]=%v, [Stack]=\n%s", v, debug.Stack())
		}
	}()

	return dfs(obj, nil, root)
}

type node struct {
	field    string
	children map[string]*node
	handler  func(obj interface{}, fields []string) error // obj可能是map[string]interface{}或[]interface{}且为key的父节点
}

// 在以n为根的trie子树中插入路径fields
func (n *node) insert(fields []string, handler func(interface{}, []string) error) error {
	fieldCopy := fields
	for {
		if len(fields) == 0 {
			if n.handler != nil {
				return fmt.Errorf("handler for '%v' is already registered", fieldCopy)
			}
			n.handler = handler
			return nil
		}
		v := n.children[fields[0]]
		if v == nil {
			v = &node{
				field:    fields[0],
				children: map[string]*node{},
			}
			n.children[fields[0]] = v
		}

		n = v
		fields = fields[1:]
	}
}

func dfs(obj interface{}, fields []string, trieRoot *node) error {
	if arr, ok := obj.([]interface{}); ok {
		for _, v := range arr {
			if err := dfs(v, fields, trieRoot); err != nil {
				return err
			}
		}
		return nil
	}

	if m, ok := obj.(map[string]interface{}); ok {
		for field, value := range m {
			n := trieRoot.children[field]
			if n == nil {
				continue
			}
			if n.handler != nil {
				if err := n.handler(obj, append(fields, field)); err != nil {
					return err
				}
			}
			if err := dfs(value, append(fields, field), n); err != nil {
				return err
			}
		}
	}
	return nil
}
