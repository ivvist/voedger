//go:build !tinygo

/*
* Copyright (c) 2023-present unTill Pro, Ltd.
*  @author Michael Saigachenko
 */

package exttinygo

import (
	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istructs"
)

type extint = int

func hostPanic(msgPtr, msgSize uint32) {

}

func hostRowWriterPutString(id uint64, typ uint32, namePtr, nameSize, valuePtr, valueSize uint32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutString(name, decodeStr(valuePtr, valueSize))
}

func hostRowWriterPutBytes(id uint64, typ uint32, namePtr, nameSize, valuePtr, valueSize uint32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)

	var bytes []byte
	val := decodeStr(valuePtr, valueSize)
	bytes = []byte(val)
	writer.PutBytes(name, bytes)
}

func hostRowWriterPutQName(id uint64, typ uint32, namePtr, nameSize, pkgPtr, pkgSize, entityPtr, entitySize uint32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	pkg := decodeStr(pkgPtr, pkgSize)
	entity := decodeStr(entityPtr, entitySize)
	localPkg := mockEngine.io.PackageLocalName(pkg)
	writer.PutQName(name, appdef.NewQName(localPkg, entity))
}

func hostRowWriterPutBool(id uint64, typ uint32, namePtr, nameSize, value uint32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutBool(name, value > 0)
}

func hostRowWriterPutInt32(id uint64, typ uint32, namePtr, nameSize, value uint32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutInt32(name, int32(value))
}

func hostRowWriterPutInt64(id uint64, typ uint32, namePtr, nameSize uint32, value uint64) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutInt64(name, int64(value))
}

func hostRowWriterPutFloat32(id uint64, typ uint32, namePtr, nameSize uint32, value float32) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutFloat32(name, value)
}

func hostRowWriterPutFloat64(id uint64, typ uint32, namePtr, nameSize uint32, value float64) {
	writer, name := mockEngine.getWriterArgs(id, typ, namePtr, nameSize)
	writer.PutFloat64(name, value)
}

func hostGetKey(storagePtr, storageSize, entityPtr, entitySize uint32) uint64 {
	var storage appdef.QName
	var entity appdef.QName
	var err error
	storage = mockEngine.parseQname(decodeStr(storagePtr, storageSize))
	entitystr := decodeStr(entityPtr, entitySize)
	if entitystr != "" {
		entity, err = appdef.ParseQName(entitystr)
		if err != nil {
			panic(err)
		}
	}
	k, e := mockEngine.io.KeyBuilder(storage, entity)
	if e != nil {
		panic(e)
	}
	res := uint64(len(mockEngine.keyBuilders))
	mockEngine.keyBuilders = append(mockEngine.keyBuilders, k)
	return res
}

func hostQueryValue(keyId uint64) (result uint64) {
	if int(keyId) >= len(mockEngine.keyBuilders) {
		panic(PanicIncorrectKeyBuilder)
	}
	v, ok, e := mockEngine.io.CanExist(mockEngine.keyBuilders[keyId])
	if e != nil {
		panic(e)
	}
	if !ok {
		return maxUint64
	}
	result = uint64(len(mockEngine.values))
	mockEngine.values = append(mockEngine.values, v)
	return
}

func hostNewValue(keyId uint64) uint64 {
	if int(keyId) >= len(mockEngine.keyBuilders) {
		panic(PanicIncorrectKeyBuilder)
	}
	vb, err := mockEngine.io.NewValue(mockEngine.keyBuilders[keyId])
	if err != nil {
		panic(err)
	}
	res := uint64(len(mockEngine.valueBuilders))
	mockEngine.valueBuilders = append(mockEngine.valueBuilders, vb)
	return res
}

func hostUpdateValue(keyId uint64, existingValueId uint64) uint64 {
	if int(keyId) >= len(mockEngine.keyBuilders) {
		panic(PanicIncorrectKeyBuilder)
	}
	if int(existingValueId) >= len(mockEngine.values) {
		panic(PanicIncorrectValue)
	}
	vb, err := mockEngine.io.UpdateValue(mockEngine.keyBuilders[keyId], mockEngine.values[existingValueId])
	if err != nil {
		panic(err)
	}
	res := uint64(len(mockEngine.valueBuilders))
	mockEngine.valueBuilders = append(mockEngine.valueBuilders, vb)
	return res
}

func hostValueLength(id uint64) uint32 {
	if int(id) >= len(mockEngine.values) {
		panic(PanicIncorrectValue)
	}
	return uint32(mockEngine.values[id].Length())
}

func hostValueAsBytes(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return mockEngine.allocAndSend(v.AsBytes(name))
}

func hostValueAsString(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return mockEngine.allocAndSend([]byte(v.AsString(name)))
}

func hostValueAsInt32(id uint64, namePtr, nameSize uint32) uint32 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return uint32(v.AsInt32(name))
}

func hostValueAsInt64(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return uint64(v.AsInt64(name))
}

func hostValueAsFloat32(id uint64, namePtr, nameSize uint32) float32 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return v.AsFloat32(name)
}

func hostValueAsFloat64(id uint64, namePtr, nameSize uint32) float64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	return v.AsFloat64(name)
}

func hostValueAsValue(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	value := v.AsValue(name)
	res := uint64(len(mockEngine.values))
	mockEngine.values = append(mockEngine.values, value)
	return res
}

func hostValueAsQNamePkg(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	qname := v.AsQName(name)
	fullPkg := mockEngine.io.PackageFullPath(qname.Pkg())
	return mockEngine.allocAndSend([]byte(fullPkg))
}

func hostValueAsQNameEntity(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	value := v.AsValue(name)
	res := uint64(len(mockEngine.values))
	mockEngine.values = append(mockEngine.values, value)
	return res
}

func hostValueAsBool(id uint64, namePtr, nameSize uint32) uint64 {
	v, name := mockEngine.valueargs(id, namePtr, nameSize)
	if v.AsBool(name) {
		return 1
	}
	return 0
}

func hostValueGetAsBytes(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	return mockEngine.allocAndSend(v.GetAsBytes(int(index)))
}

func hostValueGetAsString(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	return mockEngine.allocAndSend([]byte(v.GetAsString(int(index))))
}

func hostValueGetAsInt32(id uint64, index uint32) uint32 {
	v := mockEngine.value(id)
	return uint32(v.GetAsInt32(int(index)))
}

func hostValueGetAsInt64(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	return uint64(v.GetAsInt64(int(index)))
}

func hostValueGetAsFloat32(id uint64, index uint32) float32 {
	v := mockEngine.value(id)
	return v.GetAsFloat32(int(index))
}

func hostValueGetAsFloat64(id uint64, index uint32) float64 {
	v := mockEngine.value(id)
	return v.GetAsFloat64(int(index))
}

func hostValueGetAsValue(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	value := v.GetAsValue(int(index))
	res := uint64(len(mockEngine.values))
	mockEngine.values = append(mockEngine.values, value)
	return res
}

func hostValueGetAsQNamePkg(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	qname := v.GetAsQName(int(index))
	fullPkg := mockEngine.io.PackageFullPath(qname.Pkg())
	return mockEngine.allocAndSend([]byte(fullPkg))
}

func hostValueGetAsQNameEntity(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	qname := v.GetAsQName(int(index))
	return mockEngine.allocAndSend([]byte(qname.Entity()))
}

func hostValueGetAsBool(id uint64, index uint32) uint64 {
	v := mockEngine.value(id)
	if v.GetAsBool(int(index)) {
		return 1
	}
	return 0
}

func hostKeyAsString(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return mockEngine.allocAndSend([]byte(key.AsString(name)))
}

func hostKeyAsBytes(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return mockEngine.allocAndSend(key.AsBytes(name))
}

func hostKeyAsQNamePkg(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	qname := key.AsQName(name)
	fullPkg := mockEngine.io.PackageFullPath(qname.Pkg())
	return mockEngine.allocAndSend([]byte(fullPkg))
}

func hostKeyAsQNameEntity(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	qname := key.AsQName(name)
	return mockEngine.allocAndSend([]byte(qname.Entity()))
}

func hostKeyAsBool(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	if key.AsBool(name) {
		return uint64(1)
	}
	return uint64(0)
}

func hostKeyAsInt32(id uint64, namePtr, nameSize uint32) uint32 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return uint32(key.AsInt32(name))
}

func hostKeyAsInt64(id uint64, namePtr, nameSize uint32) uint64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return uint64(key.AsInt64(name))
}

func hostKeyAsFloat32(id uint64, namePtr, nameSize uint32) float32 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return key.AsFloat32(name)
}

func hostKeyAsFloat64(id uint64, namePtr, nameSize uint32) float64 {
	key, name := mockEngine.keyargs(id, namePtr, nameSize)
	return key.AsFloat64(name)
}

func hostReadValues(keyId uint64) {
	if int(keyId) >= len(mockEngine.keyBuilders) {
		panic(PanicIncorrectKeyBuilder)
	}
	first := true
	keyIndex := len(mockEngine.keys)
	valueIndex := len(mockEngine.values)
	err := mockEngine.io.Read(mockEngine.keyBuilders[keyId], func(key istructs.IKey, value istructs.IStateValue) (err error) {
		if first {
			mockEngine.keys = append(mockEngine.keys, key)
			mockEngine.values = append(mockEngine.values, value)
			first = false
		} else { // replace
			mockEngine.keys[keyIndex] = key
			mockEngine.values[valueIndex] = value
		}
		onReadValue(uint64(keyIndex), uint64(valueIndex))
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
}

func hostGetValue(keyId uint64) (result uint64) {
	if int(keyId) >= len(mockEngine.keyBuilders) {
		panic(PanicIncorrectKeyBuilder)
	}
	v, e := mockEngine.io.MustExist(mockEngine.keyBuilders[keyId])
	if e != nil {
		panic(e)
	}
	res := uint64(len(mockEngine.values))
	mockEngine.values = append(mockEngine.values, v)
	return res
}
