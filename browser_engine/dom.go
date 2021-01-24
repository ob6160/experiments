package main

type DOMNode struct {
  children []DOMNode
}

type DOMTextNode struct {
  DOMNode
  value string
}

type DOMElementNode struct {
  DOMNode
  tagName string
  attributes map[string]string
}

func NewDOMTextNode(value string) *DOMTextNode {
  var node = DOMNode{
    children: nil,
  }
  return &DOMTextNode{
    node,
    value,
  }
}

func NewDOMElementNode(tagName string, attributes map[string]string, children []DOMNode) *DOMElementNode {
  var node = DOMNode{
    children,
  }
  return &DOMElementNode{node, tagName, attributes}
}


