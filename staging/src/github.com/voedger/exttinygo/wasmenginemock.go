//go:build !tinygo

/*
  - Copyright (c) 2023-present unTill Software Development Group B.V.
    @author Michael Saigachenko
*/

package exttinygo

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istructs"
)

const (
	keysCapacity          = 10
	keysBuildersCapacity  = 10
	valuesCapacity        = 10
	valueBuildersCapacity = 10
	errCapacity           = 100
	maxUint64             = ^uint64(0)
)

var (
	PanicMockEngineNotDefined = "Test wasm engine not define. Call InitTestEngine first."
	PanicIncorrectKey         = "incorrect key"
	PanicIncorrectKeyBuilder  = "incorrect key builder"
	PanicIncorrectValue       = "incorrect value"
	PanicIncorrectIntent      = "incorrect intent"
	PanicIOnotassigned        = "IO not set"
)

type TState uint64

// NewKey returns a Key builder for specified storage and entity name
func (st TState) KeyBuilder(storage, entity appdef.QName) (builder istructs.IStateKeyBuilder, err error) {
	return nil, nil
}

func (st TState) CanExist(key istructs.IStateKeyBuilder) (value istructs.IStateValue, ok bool, err error) {
	return nil, false, nil
}

func (st TState) CanExistAll(keys []istructs.IStateKeyBuilder, callback istructs.StateValueCallback) (err error) {
	return nil
}

func (st TState) MustExist(key istructs.IStateKeyBuilder) (value istructs.IStateValue, err error) {
	return nil, nil
}

func (st TState) MustExistAll(keys []istructs.IStateKeyBuilder, callback istructs.StateValueCallback) (err error) {
	return nil
}

func (st TState) MustNotExist(key istructs.IStateKeyBuilder) (err error) {
	return nil
}

func (st TState) MustNotExistAll(keys []istructs.IStateKeyBuilder) (err error) {
	return nil
}

// Read reads all values according to the get and return them in callback
func (st TState) Read(key istructs.IStateKeyBuilder, callback istructs.ValueCallback) (err error) {
	return nil
}

type TExtensionIO struct {
	istructs.IState
	istructs.IIntents
	istructs.IPkgNameResolver
}

type mockExtEngine struct {
	t             *testing.T
	keys          []istructs.IKey
	keyBuilders   []istructs.IStateKeyBuilder
	values        []istructs.IStateValue
	valueBuilders []istructs.IStateValueBuilder
	errs          []string

	io TExtensionIO // mockExtension
}

func errUndefinedPackage(name string) error {
	return errors.New("undefined package: " + name)
}

var mockEngine *mockExtEngine

// InitTest - call before testing
func InitTest(t *testing.T) {
	mockEngine = &mockExtEngine{t: t}

	mockEngine.keys = make([]istructs.IKey, 0, keysCapacity)
	mockEngine.keyBuilders = make([]istructs.IStateKeyBuilder, 0, keysBuildersCapacity)
	mockEngine.values = make([]istructs.IStateValue, 0, valuesCapacity)
	mockEngine.valueBuilders = make([]istructs.IStateValueBuilder, 0, valueBuildersCapacity)
	mockEngine.errs = make([]string, 0, errCapacity)

	mockEngine.io = TExtensionIO{}
}

func (f *mockExtEngine) parseQname(value string) (qname appdef.QName) {

	pos := strings.LastIndex(value, ".")
	if pos == -1 {
		panic(fmt.Errorf("%w: %v", appdef.ErrInvalidQNameStringRepresentation, value))
	}

	packageFullPath := value[:pos]
	entityName := value[pos+1:]
	localName := f.io.PackageLocalName(packageFullPath)
	if localName == "" {
		panic(errUndefinedPackage(packageFullPath))
	}

	return appdef.NewQName(localName, entityName)
}

func (f *mockExtEngine) getWriterArgs(id uint64, typ uint32, namePtr uint32, nameSize uint32) (writer istructs.IRowWriter, name string) {
	switch typ {
	case 0:
		if int(id) >= len(f.keyBuilders) {
			panic(PanicIncorrectKeyBuilder)
		}
		writer = f.keyBuilders[id]
	default:
		if int(id) >= len(f.valueBuilders) {
			panic(PanicIncorrectIntent)
		}
		writer = f.valueBuilders[id]
	}
	name = decodeStr(namePtr, nameSize)
	return
}

func (f *mockExtEngine) allocAndSend(buf []byte) (result uint64) {
	ptr := uint32(uintptr(unsafe.Pointer(&buf)))
	return (uint64(ptr) << uint64(bitsInFourBytes)) | uint64(len(buf))
}

func (f *mockExtEngine) keyargs(id uint64, namePtr uint32, nameSize uint32) (istructs.IKey, string) {
	if int(id) >= len(f.keys) {
		panic(PanicIncorrectKey)
	}
	return f.keys[id], decodeStr(namePtr, nameSize)
}

func (f *mockExtEngine) valueargs(id uint64, namePtr uint32, nameSize uint32) (istructs.IStateValue, string) {
	if int(id) >= len(f.values) {
		panic(PanicIncorrectValue)
	}
	return f.values[id], decodeStr(namePtr, nameSize)
}

func (f *mockExtEngine) value(id uint64) istructs.IStateValue {
	if int(id) >= len(f.values) {
		panic(PanicIncorrectValue)
	}
	return f.values[id]
}

func DeInitTestEngine(engine *mockExtEngine) {
	if nil == engine {
		return
	}
	engine = nil
}

// RunAndCheck
func RunAndCheck(stateIntents []TIntent, testFunc func(...any), expectedIntents []TIntent) {
	var resultIntents []TIntent

	require.Equal(mockEngine.t, len(expectedIntents), len(resultIntents))
}
