package main

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  children []*DOMNode
}

//
//func NewDOMTextNode(value string) *DOMTextNode {
//  return &DOMTextNode{
//    value,
//  }
//}
//
//func NewDOMElementNode(tagName string, attributes map[string]string, children []DOMNode) *DOMElementNode {
//  return &DOMElementNode{children, tagName, attributes}
//}


