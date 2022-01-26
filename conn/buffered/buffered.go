/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

// Package buffered creates buffered net.Conn connections
// that will continuously read in a separate goroutine,
// allowing for testing packages to write without needing
// to concurrently read
package buffered

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

var _ net.Conn = (*Buffered)(nil)

// Buffered is a buffered net.Conn that will continuously
// read from the underlying connection in a separate goroutine
type Buffered struct {
	c    net.Conn
	wg   sync.WaitGroup
	buf  *bytes.Buffer
	cond *sync.Cond
	err  error
	size int
}

// New creates a new Buffered net.Conn with a specified read size.
// This size defines the maximum amount of data that will be read
// at a time in the background.
func New(c net.Conn, size int) (a *Buffered) {
	a = &Buffered{
		c:    c,
		buf:  new(bytes.Buffer),
		cond: sync.NewCond(new(sync.Mutex)),
		size: size,
	}

	a.wg.Add(1)
	go a.readLoop()

	return a
}

func (a *Buffered) readLoop() {
	defer a.wg.Done()
	data := make([]byte, a.size)
LOOP:
	n, err := a.c.Read(data[0:])
	if err != nil {
		a.cond.L.Lock()
		a.err = err
		a.cond.L.Unlock()
		a.cond.Signal()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.ErrClosedPipe) {
			return
		}
		goto LOOP
	}
	a.cond.L.Lock()
	a.err = nil
	a.buf.Write(data[:n])
	a.cond.L.Unlock()
	a.cond.Signal()
	goto LOOP
}

// Read satisfies the io.Reader interface and allow safe
// concurrent reads from the buffered connection
func (a *Buffered) Read(b []byte) (int, error) {
	a.cond.L.Lock()
	defer a.cond.L.Unlock()
LOOP:
	if a.err != nil {
		return 0, a.err
	}
	if a.buf.Len() == 0 && len(b) > 0 {
		a.cond.Wait()
		goto LOOP
	}
	return a.buf.Read(b)
}

// Write satisfies the io.Writer interface and
// transparently writes to the underlying net.Conn
func (a *Buffered) Write(b []byte) (int, error) {
	return a.c.Write(b)
}

// Close satisfies the io.Closer interface and
// closes the underlying net.Conn as well as
// stopping the background read goroutine
func (a *Buffered) Close() error {
	defer a.wg.Wait()
	return a.c.Close()
}

// LocalAddr helps satisfy the net.Conn interface
// and returns the local address of the underlying net.Conn
func (a *Buffered) LocalAddr() net.Addr {
	return a.c.LocalAddr()
}

// RemoteAddr helps satisfy the net.Conn interface
// and returns the remote address of the underlying net.Conn
func (a *Buffered) RemoteAddr() net.Addr {
	return a.c.RemoteAddr()
}

// SetDeadline helps satisfy the net.Conn interface
// and sets a deadline on the underlying net.Conn
func (a *Buffered) SetDeadline(t time.Time) error {
	return a.c.SetDeadline(t)
}

// SetReadDeadline helps satisfy the net.Conn interface
// and sets a read deadline on the underlying net.Conn
func (a *Buffered) SetReadDeadline(t time.Time) error {
	return a.c.SetReadDeadline(t)
}

// SetWriteDeadline helps satisfy the net.Conn interface
// and sets a write deadline on the underlying net.Conn
func (a *Buffered) SetWriteDeadline(t time.Time) error {
	return a.c.SetWriteDeadline(t)
}
