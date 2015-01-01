package store

// #include <rocksdb/c.h>
// #include <malloc.h>
// #include <string.h>
import "C"

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"unsafe"
)

const (
	TotalUnitNum = 4096
)

const (
	ColSpaceDefault     = 0x1
	ColSpaceScore       = 0x2
	ColSpaceScoreSorted = 0x4
)

type TableKey struct {
	TableId  uint8
	RowKey   []byte // len(rowKey) < 256
	ColSpace uint8
	ColKey   []byte
}

func GetUnitId(rowKey []byte) uint16 {
	return uint16(crc32.ChecksumIEEE(rowKey) % TotalUnitNum)
}

func GetRawKey(dbId uint8, key TableKey) []byte {
	var unitId = GetUnitId(key.RowKey)

	// wUnitId+cDbId+cTableId+cKeyLen+sRowKey+colType+sColKey
	var rawLen = 6 + len(key.RowKey) + len(key.ColKey)
	var rawKey = make([]byte, rawLen, rawLen)
	binary.BigEndian.PutUint16(rawKey, unitId)
	rawKey[2] = dbId
	rawKey[3] = key.TableId
	rawKey[4] = uint8(len(key.RowKey))
	copy(rawKey[5:], key.RowKey)
	rawKey[5+len(key.RowKey)] = key.ColSpace
	copy(rawKey[(6+len(key.RowKey)):], key.ColKey)

	return rawKey
}

func ParseRawKey(rawKey []byte) (unitId uint16, dbId uint8, key TableKey) {
	unitId = binary.BigEndian.Uint16(rawKey)
	dbId = rawKey[2]
	key.TableId = rawKey[3]
	var keyLen = rawKey[4]
	var colTypePos = 5 + int(keyLen)
	key.RowKey = rawKey[5:colTypePos]
	key.ColSpace = rawKey[colTypePos]
	key.ColKey = rawKey[(colTypePos + 1):]
	return
}

type TableDB struct {
	db   *C.rocksdb_t
	opt  *C.rocksdb_options_t
	rOpt *C.rocksdb_readoptions_t
	wOpt *C.rocksdb_writeoptions_t
}

type TableIterator struct {
	iter *C.rocksdb_iterator_t
}

func NewTableDB() *TableDB {
	db := new(TableDB)

	return db
}

func (db *TableDB) Close() {
	if db.db != nil {
		C.rocksdb_close(db.db)
		db.db = nil

		C.rocksdb_options_destroy(db.opt)
		C.rocksdb_readoptions_destroy(db.rOpt)
		C.rocksdb_writeoptions_destroy(db.wOpt)
	}
}

func (db *TableDB) Open(name string, createIfMissing bool) error {
	var errStr *C.char

	db.opt = C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(db.opt,
		boolToUchar(createIfMissing))
	C.rocksdb_options_set_write_buffer_size(db.opt, 1024*1024*64)

	var block_cache = C.rocksdb_cache_create_lru(1024 * 1024 * 64)
	var block_cache_compressed = C.rocksdb_cache_create_lru(1024 * 1024 * 64)
	var block_based_table_options = C.rocksdb_block_based_options_create()
	C.rocksdb_block_based_options_set_block_cache_compressed(
		block_based_table_options, block_cache_compressed)
	C.rocksdb_block_based_options_set_block_cache(
		block_based_table_options, block_cache)
	C.rocksdb_options_set_block_based_table_factory(
		db.opt, block_based_table_options)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	db.db = C.rocksdb_open(db.opt, cname, &errStr)
	if errStr != nil {
		defer C.free(unsafe.Pointer(errStr))
		return fmt.Errorf(C.GoString(errStr))
	}

	db.rOpt = C.rocksdb_readoptions_create()
	db.wOpt = C.rocksdb_writeoptions_create()

	return nil
}

func (db *TableDB) Put(dbId uint8, key TableKey, value []byte) error {
	var rawKey = GetRawKey(dbId, key)

	var ck, cv *C.char
	if len(rawKey) > 0 {
		ck = (*C.char)(unsafe.Pointer(&rawKey[0]))
	}
	if len(value) > 0 {
		cv = (*C.char)(unsafe.Pointer(&value[0]))
	}

	var errStr *C.char
	C.rocksdb_put(db.db, db.wOpt, ck, C.size_t(len(rawKey)), cv, C.size_t(len(value)),
		&errStr)
	if errStr != nil {
		defer C.free(unsafe.Pointer(errStr))
		return fmt.Errorf(C.GoString(errStr))
	}

	return nil
}

func (db *TableDB) Get(dbId uint8, key TableKey) ([]byte, error) {
	var rawKey = GetRawKey(dbId, key)
	var ck = (*C.char)(unsafe.Pointer(&rawKey[0]))

	var errStr *C.char
	var vallen C.size_t
	var cv = C.rocksdb_get(db.db, db.rOpt, ck, C.size_t(len(rawKey)), &vallen, &errStr)

	var err error
	if errStr != nil {
		defer C.free(unsafe.Pointer(errStr))
		err = fmt.Errorf(C.GoString(errStr))
	}

	if cv != nil {
		defer C.free(unsafe.Pointer(cv))
		return C.GoBytes(unsafe.Pointer(cv), C.int(vallen)), err
	}

	return nil, err
}

func (db *TableDB) Mput(dbId uint8, keys []TableKey, values [][]byte) error {
	if len(keys) != len(values) {
		return fmt.Errorf("invalid keys or values")
	}

	if len(keys) == 0 {
		return nil // nothing to do
	}

	var batch = C.rocksdb_writebatch_create()
	defer C.rocksdb_writebatch_destroy(batch)

	for i := 0; i < len(keys); i++ {
		var value = values[i]
		var rawKey = GetRawKey(dbId, keys[i])

		var ck = (*C.char)(unsafe.Pointer(&rawKey[0]))
		var cv *C.char
		if len(value) > 0 {
			cv = (*C.char)(unsafe.Pointer(&value[0]))
		}

		C.rocksdb_writebatch_put(batch, ck, C.size_t(len(rawKey)), cv, C.size_t(len(value)))
	}

	var errStr *C.char

	C.rocksdb_write(db.db, db.wOpt, batch, &errStr)
	if errStr != nil {
		defer C.free(unsafe.Pointer(errStr))
		return fmt.Errorf(C.GoString(errStr))
	}

	return nil
}

func (db *TableDB) NewIterator(fillCache bool) *TableIterator {
	var iter = new(TableIterator)
	var scanOpt = C.rocksdb_readoptions_create()
	defer C.rocksdb_readoptions_destroy(scanOpt)

	C.rocksdb_readoptions_set_fill_cache(scanOpt, boolToUchar(fillCache))
	iter.iter = C.rocksdb_create_iterator(db.db, scanOpt)

	return iter
}

func (iter *TableIterator) Close() {
	if iter.iter != nil {
		C.rocksdb_iter_destroy(iter.iter)
		iter.iter = nil
	}
}

func (iter *TableIterator) SeekToFirst() {
	C.rocksdb_iter_seek_to_first(iter.iter)
}

func (iter *TableIterator) SeekToLast() {
	C.rocksdb_iter_seek_to_last(iter.iter)
}

func (iter *TableIterator) Seek(key []byte) {
	var ck *C.char
	if len(key) > 0 {
		ck = (*C.char)(unsafe.Pointer(&key[0]))
	}
	C.rocksdb_iter_seek(iter.iter, ck, C.size_t(len(key)))
}

func (iter *TableIterator) Next() {
	C.rocksdb_iter_next(iter.iter)
}

func (iter *TableIterator) Prev() {
	C.rocksdb_iter_prev(iter.iter)
}

func (iter *TableIterator) Valid() bool {
	return C.rocksdb_iter_valid(iter.iter) != 0
}

func (iter *TableIterator) Key() (unitId uint16, dbId uint8, key TableKey) {
	var keyLen C.size_t
	var ck = C.rocksdb_iter_key(iter.iter, &keyLen)
	var rawKey = C.GoBytes(unsafe.Pointer(ck), C.int(keyLen))
	return ParseRawKey(rawKey)
}

func (iter *TableIterator) Value() []byte {
	var valueLen C.size_t
	var value = C.rocksdb_iter_value(iter.iter, &valueLen)
	return C.GoBytes(unsafe.Pointer(value), C.int(valueLen))
}

func boolToUchar(b bool) C.uchar {
	if b {
		return 1
	} else {
		return 0
	}
}
