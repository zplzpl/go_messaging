package unique_id

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	ErrEmptyUniqueID = errors.New("unique id is empty")
)

var node *snowflake.Node
var once sync.Once

type UniqueID snowflake.ID

func nodeID() int64 {
	s := rand.NewSource(time.Now().UnixNano())
	return s.Int63() % 1024
}

func GenerateID() UniqueID {
	once.Do(func() {
		var err error
		node, err = snowflake.NewNode(nodeID())
		if err != nil {
			panic(err)
		}
	})
	id := node.Generate()
	return UniqueID(id)
}

func (v UniqueID) Int64() int64 {
	return int64(v)
}

func (v UniqueID) UInt64() uint64 {
	return uint64(v)
}

func (v UniqueID) Equal(id UniqueID) bool {
	return v == id
}

func (v UniqueID) Validate() error {

	if v.IsEmpty() {
		return ErrEmptyUniqueID
	}
	return nil
}

func (v UniqueID) IsEmpty() bool {
	if v == 0 {
		return true
	}
	return false
}

func ParseUniqueID(id uint64) UniqueID {
	return UniqueID(id)
}
