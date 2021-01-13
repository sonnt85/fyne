package repository

import (
	"io/ioutil"
	"testing"

	"fyne.io/fyne/storage"
	"fyne.io/fyne/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRepositoryRegistration(t *testing.T) {
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)

	// this should never fail, and we assume it doesn't in other tests here
	// for brevity
	foo, err := storage.ParseURI("mem://foo")
	assert.Nil(t, err)

	// make sure we get the same repo back
	repo, err := repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, m, repo)

	// test that re-registration also works
	m2 := NewMemoryRepository("mem")
	repository.Register("mem", m2)
	assert.False(t, m == m2) // this is explicitly intended to be pointer comparison
	repo, err = repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, m2, repo)
}

func TestMemoryRepositoryExists(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo"] = []byte{}
	m.data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooExists, err := storage.Exists(foo)
	assert.True(t, fooExists)
	assert.Nil(t, err)

	barExists, err := storage.Exists(bar)
	assert.True(t, barExists)
	assert.Nil(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.Nil(t, err)
}

func TestMemoryRepositoryReader(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo"] = []byte{}
	m.data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{}, fooData)
	assert.Nil(t, err)

	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{1, 2, 3}, barData)
	assert.Nil(t, err)

	bazReader, err := storage.Reader(baz)
	assert.Nil(t, bazReader)
	assert.NotNil(t, err)
}

func TestMemoryRepositoryCanRead(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo"] = []byte{}
	m.data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooCanRead, err := storage.CanRead(foo)
	assert.True(t, fooCanRead)
	assert.Nil(t, err)

	barCanRead, err := storage.CanRead(bar)
	assert.True(t, barCanRead)
	assert.Nil(t, err)

	bazCanRead, err := storage.CanRead(baz)
	assert.False(t, bazCanRead)
	assert.NotNil(t, err)
}

func TestMemoryRepositoryWriter(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo"] = []byte{}
	m.data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	// write some data and assert there are no errors
	fooWriter, err := storage.Writer(foo)
	assert.Nil(t, err)
	assert.NotNil(t, fooWriter)

	barWriter, err := storage.Writer(bar)
	assert.Nil(t, err)
	assert.NotNil(t, barWriter)

	bazWriter, err := storage.Writer(baz)
	assert.Nil(t, err)
	assert.NotNil(t, bazWriter)

	n, err := fooWriter.Write([]byte{1, 2, 3, 4, 5})
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = barWriter.Write([]byte{6, 7, 8, 9})
	assert.Nil(t, err)
	assert.Equal(t, 4, n)

	n, err = bazWriter.Write([]byte{5, 4, 3, 2, 1, 0})
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	fooWriter.Close()
	barWriter.Close()
	bazWriter.Close()

	// now make sure we can read the data back correctly
	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, fooData)
	assert.Nil(t, err)

	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{6, 7, 8, 9}, barData)
	assert.Nil(t, err)

	bazReader, err := storage.Reader(baz)
	assert.Nil(t, err)
	bazData, err := ioutil.ReadAll(bazReader)
	assert.Equal(t, []byte{5, 4, 3, 2, 1, 0}, bazData)
	assert.Nil(t, err)

	// now let's test deletion
	err = storage.Delete(foo)
	assert.Nil(t, err)

	err = storage.Delete(bar)
	assert.Nil(t, err)

	err = storage.Delete(baz)
	assert.Nil(t, err)

	fooExists, err := storage.Exists(foo)
	assert.False(t, fooExists)
	assert.Nil(t, err)

	barExists, err := storage.Exists(bar)
	assert.False(t, barExists)
	assert.Nil(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.Nil(t, err)

}

func TestMemoryRepositoryCanWrite(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo"] = []byte{}
	m.data["/bar"] = []byte{1, 2, 3}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo")
	bar, _ := storage.ParseURI("mem:///bar")
	baz, _ := storage.ParseURI("mem:///baz")

	fooCanWrite, err := storage.CanWrite(foo)
	assert.True(t, fooCanWrite)
	assert.Nil(t, err)

	barCanWrite, err := storage.CanWrite(bar)
	assert.True(t, barCanWrite)
	assert.Nil(t, err)

	bazCanWrite, err := storage.CanWrite(baz)
	assert.True(t, bazCanWrite)
	assert.Nil(t, err)
}

func TestMemoryRepositoryParent(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)
	m.data["/foo/bar/baz"] = []byte{}

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo/bar/baz")
	fooExpectedParent, _ := storage.ParseURI("mem:///foo/bar")
	fooExists, err := storage.Exists(foo)
	assert.True(t, fooExists)
	assert.Nil(t, err)

	fooParent, err := storage.Parent(foo)
	assert.Nil(t, err)
	assert.Equal(t, fooExpectedParent.String(), fooParent.String())
}

func TestMemoryRepositoryChild(t *testing.T) {
	// set up our repository - it's OK if we already registered it
	m := NewMemoryRepository("mem")
	repository.Register("mem", m)

	// and some URIs - we know that they will not fail parsing
	foo, _ := storage.ParseURI("mem:///foo/bar/baz")
	fooExpectedChild, _ := storage.ParseURI("mem:///foo/bar/baz/quux")

	fooChild, err := storage.Child(foo, "quux")
	assert.Nil(t, err)
	assert.Equal(t, fooExpectedChild.String(), fooChild.String())
}