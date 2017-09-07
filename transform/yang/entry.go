package yang

import (
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/oshothebig/pbast"
)

type entry struct {
	*yang.Entry
}

func (e entry) rpcs() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.RPC != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	rpcs := []entry{}
	for _, name := range names {
		rpcs = append(rpcs, entry{e.Dir[name]})
	}

	return rpcs
}

func (e entry) notifications() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.Kind == yang.NotificationEntry {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	ns := []entry{}
	for _, name := range names {
		ns = append(ns, entry{e.Dir[name]})
	}

	return ns
}

func (e entry) children() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.RPC != nil || child.Kind == yang.NotificationEntry {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	children := []entry{}
	for _, name := range names {
		children = append(children, entry{e.Dir[name]})
	}

	return children
}

func (e entry) moduleComment() pbast.Comment {
	return joinComments("",
		e.description(),
		e.attributeToComment("yang-version", "YANG version"),
		e.namespace(),
		e.attributeToComment("organization", "Organization"),
		e.attributeToComment("contact", "Contact"),
		e.revisions(),
		e.reference(),
	)
}

func (e entry) genericComments() pbast.Comment {
	return joinComments("",
		e.description(),
		e.reference(),
	)
}

func (e entry) description() pbast.Comment {
	description := e.Description
	if e.Description == "" {
		return nil
	}

	lines := strings.Split(strings.TrimRight(description, "\n "), "\n")

	ret := make([]string, 0, len(lines)+1)
	ret = append(ret, "Description:")
	ret = append(ret, lines...)
	return ret
}

func (e entry) revisions() pbast.Comment {
	var lines []string
	if v := e.Extra["revision"]; len(v) > 0 {
		for _, rev := range v[0].([]*yang.Revision) {
			lines = append(lines, "Revision: "+rev.Name)
		}
	}

	return lines
}

func (e entry) attributeToComment(stmt, comment string) pbast.Comment {
	v := e.Extra[stmt]
	if len(v) == 0 {
		return nil
	}

	attribute := v[0].(*yang.Value)
	if attribute == nil {
		return nil
	}
	if attribute.Name == "" {
		return nil
	}

	return []string{comment + ": " + attribute.Name}
}

func (e entry) namespace() pbast.Comment {
	namespace := e.Namespace().Name
	if namespace == "" {
		return nil
	}

	return []string{"Namespace: " + namespace}
}

func (e entry) reference() pbast.Comment {
	v := e.Extra["reference"]
	if len(v) == 0 {
		return nil
	}

	ref := v[0].(*yang.Value)
	if ref == nil {
		return nil
	}
	if ref.Name == "" {
		return nil
	}

	lines := strings.Split(strings.TrimRight(ref.Name, "\n "), "\n")

	ret := make([]string, 0, len(lines)+1)
	ret = append(ret, "Reference:")
	ret = append(ret, lines...)
	return ret
}

func joinComments(sep string, comments ...pbast.Comment) pbast.Comment {
	if len(comments) == 0 {
		return nil
	}

	var concat []string
	concat = append(concat, comments[0]...)
	for _, comment := range comments[1:] {
		if len(comment) == 0 {
			continue
		}

		concat = append(concat, sep)
		concat = append(concat, comment...)
	}

	return concat
}
