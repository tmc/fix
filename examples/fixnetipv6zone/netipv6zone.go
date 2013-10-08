// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/tmc/fix"
	"go/ast"
)

func init() {
	fix.Register(netipv6zoneFix)
}

var netipv6zoneFix = fix.Fix{
	"netipv6zone",
	"2012-11-26",
	netipv6zone,
	`Adapt element key to IPAddr, UDPAddr or TCPAddr composite literals.

https://codereview.appspot.com/6849045/
`,
}

func netipv6zone(f *ast.File) bool {
	if !fix.Imports(f, "net") {
		return false
	}

	fixed := false
	fix.Walk(f, func(n interface{}) {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}
		se, ok := cl.Type.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if !fix.IsTopName(se.X, "net") || se.Sel == nil {
			return
		}
		switch ss := se.Sel.String(); ss {
		case "IPAddr", "UDPAddr", "TCPAddr":
			for i, e := range cl.Elts {
				if _, ok := e.(*ast.KeyValueExpr); ok {
					break
				}
				switch i {
				case 0:
					cl.Elts[i] = &ast.KeyValueExpr{
						Key:   ast.NewIdent("IP"),
						Value: e,
					}
				case 1:
					if elit, ok := e.(*ast.BasicLit); ok && elit.Value == "0" {
						cl.Elts = append(cl.Elts[:i], cl.Elts[i+1:]...)
					} else {
						cl.Elts[i] = &ast.KeyValueExpr{
							Key:   ast.NewIdent("Port"),
							Value: e,
						}
					}
				}
				fixed = true
			}
		}
	})
	return fixed
}
