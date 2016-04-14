package main

import (
	"fmt"
	"log"
	"strings"
)

type MessageElement struct {
	XSIType     string `xml:"xsitype,attr"`
	XMIId       string `xml:"http://www.omg.org/XMI id,attr"`
	Definition  string `xml:"definition,attr"`
	Name        string `xml:"name,attr"`
	MaxOccurs   string `xml:"maxOccurs,attr"`
	MinOccurs   string `xml:"minOccurs,attr"`
	XMLTag      string `xml:"xmlTag,attr"`
	SimpleType  string `xml:"simpleType,attr"`
	ComplexType string `xml:"complexType,attr"`
	Type        string `xml:"type,attr"`
}

func (m *MessageElement) IsArray() bool {
	return m.MaxOccurs != "" && m.MaxOccurs != "0" && m.MaxOccurs != "1"
}

func (m *MessageElement) ArrayDeclaration() string {
	if m.IsArray() {
		return "[]"
	}
	return ""
}

func (m *MessageElement) Declaration() string {
	return fmt.Sprintf(
		"// %s\n\t%s %s*%s `xml:\"%s%s\"`",
		strings.Replace(m.Definition, "\n", "\n\t// ", -1),
		m.Name, m.ArrayDeclaration(), m.MemberType(), m.XMLTag, m.optional())
}

type context struct {
	Receiver     string
	ReceiverType string
	Element      string
	ElementType  string
}

func (m *MessageElement) Access(receiverType string) string {
	c := &context{
		Receiver:     strings.ToLower(receiverType[:1]),
		ReceiverType: receiverType,
		Element:      m.Name,
	}
	if len(m.Type) > 0 {
		c.ElementType = typeMap[m.Type].Name
		if m.IsArray() {
			return complexArrayAccess(c)
		}
		return complexAccess(c)
	}
	if len(m.ComplexType) > 0 {
		c.ElementType = typeMap[m.ComplexType].Name
		if m.IsArray() {
			return complexArrayAccess(c)
		}
		return complexAccess(c)
	}
	if len(m.SimpleType) > 0 {
		c.ElementType = typeMap[m.SimpleType].Name
		if m.IsArray() {
			return simpleArrayAccess(c)
		}
		return simpleAccess(c)
	}
	log.Fatalf("message element with undefined type: %+v", m)
	return ""
}

func complexArrayAccess(c *context) string {
	return fmt.Sprintf("func (%s *%s) Add%s() *%s {\n\tnewValue := new(%s)\n\t%s.%s = append(%s.%s, newValue)\n\treturn newValue\n}\n",
		c.Receiver, c.ReceiverType, c.Element, c.ElementType,
		c.ElementType,
		c.Receiver, c.Element, c.Receiver, c.Element)
}

func complexAccess(c *context) string {
	return fmt.Sprintf("func (%s *%s) Add%s() *%s {\n\t%s.%s = new(%s)\n\treturn %s.%s\n}\n",
		c.Receiver, c.ReceiverType, c.Element, c.ElementType,
		c.Receiver, c.Element, c.ElementType,
		c.Receiver, c.Element)
}

func simpleArrayAccess(c *context) string {
	return fmt.Sprintf("func (%s *%s) Add%s(value string) {\n\t%s.%s = append(%s.%s, (*%s)(&value))\n}\n",
		c.Receiver, c.ReceiverType, c.Element,
		c.Receiver, c.Element, c.Receiver, c.Element, c.ElementType)
}

func simpleAccess(c *context) string {
	return fmt.Sprintf("func (%s *%s) Set%s(value string) {\n\t%s.%s = (*%s)(&value)\n}\n",
		c.Receiver, c.ReceiverType, c.Element,
		c.Receiver, c.Element, c.ElementType)
}

func (m *MessageElement) optional() string {
	if m.MinOccurs == "0" || m.MinOccurs == "" {
		return ",omitempty"
	}
	return ""
}

func (m *MessageElement) MemberType() string {
	if len(m.SimpleType) > 0 {
		return typeMap[m.SimpleType].Name
	}
	if len(m.ComplexType) > 0 {
		return typeMap[m.ComplexType].Name
	}
	if len(m.Type) > 0 {
		return typeMap[m.Type].Name
	}
	log.Fatalf("message element with undefined type: %+v", m)
	return ""
}
