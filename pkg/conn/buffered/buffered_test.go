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

package buffered

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	c1, c2 := net.Pipe()
	data := make([]byte, 2048)
	read := make([]byte, 512)
	rand.Read(data)
	var a1, a2 *Buffered

	t.Run("goroutine start", func(t *testing.T) {
		startRoutine := runtime.NumGoroutine()
		a1 = New(c1, 512)
		assert.Equal(t, startRoutine+1, runtime.NumGoroutine())
		a2 = New(c2, 512)
		assert.Equal(t, startRoutine+2, runtime.NumGoroutine())
	})

	t.Run("buffered write", func(t *testing.T) {
		n, err := a1.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
	})

	t.Run("buffered read", func(t *testing.T) {
		n, err := a2.Read(read)
		assert.NoError(t, err)
		assert.Equal(t, len(read), n)
		assert.Equal(t, data[:512], read)

		n, err = a2.Read(read[0:])
		assert.NoError(t, err)
		assert.Equal(t, len(read), n)
		assert.Equal(t, data[512:1024], read)

		n, err = a2.Read(read[0:])
		assert.NoError(t, err)
		assert.Equal(t, len(read), n)
		assert.Equal(t, data[1024:1536], read)

		n, err = a2.Read(read[0:])
		assert.NoError(t, err)
		assert.Equal(t, len(read), n)
		assert.Equal(t, data[1536:2048], read)
	})

	t.Run("shorter read", func(t *testing.T) {
		n, err := a1.Write(data[:8])
		assert.NoError(t, err)
		assert.Equal(t, 8, n)

		n, err = a2.Read(read[0:])
		assert.NoError(t, err)
		assert.Equal(t, 8, n)
		assert.Equal(t, data[:8], read[:8])
	})

	t.Run("shorter longer", func(t *testing.T) {
		n, err := a1.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)

		n, err = a2.Read(read[0:4])
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, data[:4], read[:4])
	})

	t.Run("goroutine stop", func(t *testing.T) {
		endRoutine := runtime.NumGoroutine()
		assert.NoError(t, a1.Close())
		assert.Equal(t, endRoutine-1, runtime.NumGoroutine())

		assert.NoError(t, a2.Close())
		assert.Equal(t, endRoutine-2, runtime.NumGoroutine())
	})
}
