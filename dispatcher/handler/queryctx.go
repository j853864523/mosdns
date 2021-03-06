//     Copyright (C) 2020, IrineSistiana
//
//     This file is part of mosdns.
//
//     mosdns is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     mosdns is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

package handler

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
)

type Context struct {
	Q    *dns.Msg
	From net.Addr

	R *dns.Msg
}

func (ctx *Context) Copy() *Context {
	newCtx := new(Context)
	if ctx.Q != nil {
		newCtx.Q = ctx.Q.Copy()
	}
	if ctx.R != nil {
		newCtx.R = ctx.R.Copy()
	}
	newCtx.From = ctx.From

	return newCtx
}

func (ctx *Context) String() string {
	var question []dns.Question
	if ctx.Q != nil {
		question = ctx.Q.Question
	}

	return fmt.Sprintf("%v, from: %v", question, ctx.From)
}
