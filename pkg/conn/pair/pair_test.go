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

package pair

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	data := make([]byte, 512)
	read := make([]byte, 512)
	var n int
	var err error

	n, err = rand.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	c1, c2, err := New()
	assert.NoError(t, err)

	t.Run("test write simplex", func(t *testing.T) {
		n, err = c1.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
	})

	t.Run("test read simplex", func(t *testing.T) {
		n, err = c2.Read(read)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
		assert.Equal(t, data, read)
	})

	rand.Read(data[0:])

	t.Run("test duplex", func(t *testing.T) {
		startRead := make(chan struct{}, 1)
		startWrite := make(chan struct{}, 1)
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			startRead <- struct{}{}
			read := make([]byte, 512)
			n, err := c2.Read(read)
			assert.NoError(t, err)
			assert.Equal(t, len(data), n)
			assert.Equal(t, data, read)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			startWrite <- struct{}{}
			n, err := c1.Write(data)
			assert.NoError(t, err)
			assert.Equal(t, len(data), n)
		}()

		<-startRead
		n, err = c2.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)

		<-startWrite
		n, err = c1.Read(read[0:])
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
		assert.Equal(t, data, read)

		wg.Wait()
	})

	assert.NoError(t, Cleanup(c1, c2))
}
