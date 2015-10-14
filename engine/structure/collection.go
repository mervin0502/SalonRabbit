// Package node provides a general hash based map for any type that implements
// *engine.ObjectIndex, it is ported from jdk
package structure

//import "fmt"
import (
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/node"
)
import (
	// "math"
	"math/rand"
	"time"
)

const (
	// The default initial capacity - MUST be a power of two.
	defaultInitialCapacity = 16

	// The maximum capacity, used if a higher value is implicitly specified
	// by either of the constructors with arguments.
	// MUST be a power of two <= 1<<30.
	maximumCapacity uint32 = 1 << 31

	// The load factor used when none specified in constructor.
	defaultLoadFactor = 0.75
)

type entry struct {
	// hashCode cached the hash code of key
	hashCode uint32
	// key
	key [16]byte

	// the mapped value of key
	value *node.Node

	// next entry in entry list
	next *entry
}
type entrys struct {

	//first entry
	head *entry

	//
	size uint32
}

type Collection struct {
	// The table, resized as necessary. Length MUST Always be a power of two.
	table []*entrys

	// The number of key-value mappings contained in this map.
	size uint32
	//capacity
	capacity uint32

	// The next size value at which to resize (capacity * load factor).
	threshold uint32

	// The load factor for the hash table.
	loadFactor float64

	//mutex
	//mutex *sync.RWMutex
}

// New creates hashmap with default settings.
func NewCollection() *Collection {
	c := new(Collection)
	c.loadFactor = defaultLoadFactor
	c.threshold = uint32(defaultInitialCapacity * defaultLoadFactor)
	c.capacity = defaultInitialCapacity
	c.table = make([]*entrys, defaultInitialCapacity)
	//c.mutex = new(sync.RWMutex)
	return c
}

// Applies a supplemental hash function to a given hashCode, which
// defends against poor quality hash functions.  This is critical
// because Collection uses power-of-two length hash tables, that
// otherwise encounter collisions for hashCodes that do not differ
// in lower bits. Note: Null keys always map to hash 0, thus index 0.
//
func hash(h uint32) uint32 {
	// This function ensures that hashCodes that differ only by
	// constant multiples at each bit position have a bounded
	// number of collisions (approximately 8 at default load factor).
	h ^= (h >> 20) ^ (h >> 12)
	return h ^ (h >> 7) ^ (h >> 4)
}

// Returns index for hash code c.
func indexFor(h, length uint32) uint32 {
	return h & (length - 1)
}

// Returns the number of key-value mappings in this map.
func (c Collection) Size() uint32 {
	//c.mutex.RLock()
	//defer //c.mutex.RUnlock()
	return c.size
}

// Returns true if this map contains no key-value mappings.
func (c Collection) IsEmpty() bool {
	//c.mutex.RLock()
	//defer //c.mutex.RUnlock()
	return c.size == 0
}

// Returns the find result & value to which the specified key is mapped,
// or nil if this map contains no mapping for the key.
//
// More formally, if this map contains a mapping from a key
// to a value such that (key==null ? k==null :
// key.equals(k)), then this method returns v; otherwise
// it returns nil.  (There can be at most one such mapping.)
func (c Collection) Get(key *engine.ObjectIndex) (*node.Node, bool) {
	//if key == nil {
	//	return c.getForNullKey()
	//}
	h := hash(key.GetId())
	es := c.table[indexFor(h, c.capacity)]
	if es == nil {
		return nil, false
	}
	for e := es.head; e != nil; e = e.next {
		if e.hashCode == h && e.key == key.LocalIndex() {
			return e.value, true
		}

	}
	return nil, false
}

//GetById
func (c Collection) GetById(id *engine.ObjectId) (*node.Node, bool) {

	hash := hash(id.Hash())
	es := c.table[indexFor(hash, c.capacity)]
	if es == nil {
		return nil, false
	}
	for e := es.head; e != nil; e = e.next {
		if e.value.Id.Equals(id) {
			return e.value, true
		}
	}
	return nil, false
}

//GetFirst
func (c Collection) GetFirst() (*node.Node, bool) {
	//c.mutex.RLock()
	//defer //c.mutex.RUnlock()
	var i uint32 = 0
	t := c.table[i]
	for {
		if t != nil {
			return t.head.value, true
		}
		i++
		t = c.table[i]
	}
	return nil, false
}

//GetLast
func (c Collection) GetLast() (*node.Node, bool) {
	//c.mutex.RLock()
	//defer //c.mutex.RUnlock()
	var i uint32 = 1
	t := c.table[c.capacity-i]
	for t == nil {
		i++
		t = c.table[c.capacity-i]
	}
	e := t.head
	for {
		if e.next == nil {
			return e.value, true
		}
		e = e.next
	}
	return nil, false
}

// Offloaded version of get() to look up nil keys.  Null keys map
// to index 0.  This nil case is split out uint32o separate methods
// for the sake of performance in the two most commonly used
// operations (get and put), but incorporated with conditionals in
// otherst.
/*
func (c Collection) getForNullKey() (*Node, bool) {
	for e := c.table[0].head; e != nil; e = e.next {
		if e.key == nil {
			return e.value, true
		}
	}
	return nil, false
}
*/

// Returns true if this map contains a mapping for the
// specified key.
func (c Collection) ContainsKey(key *engine.ObjectIndex) bool {
	return c.getEntry(key) != nil
}

// Returns the entry associated with the specified key in the
// Collection.  Returns nil if the Collection contains no mapping
// for the key.
func (c Collection) getEntry(key *engine.ObjectIndex) *entry {
	//hashCode := 0
	//if key != nil {
	//	hashCode = hash(key.HashCode())
	//}
	h := hash(key.GetId())
	es := c.table[indexFor(h, c.capacity)]
	if es == nil {
		return nil
	}
	for e := es.head; e != nil; e = e.next {
		if e.hashCode == key.LocalHash() && e.key == key.LocalIndex() {
			return e
		}

	}
	return nil
}

// Associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key, the old
// value is replaced.
// return a bool indicate if there already has a mapping for the key
// and the previous value associated with key, or nil if there was no mapping for key.

func (c *Collection) Add(id *engine.ObjectId) (*node.Node, bool) {
	print(id.Value, "\t")
	println(c.size, 1111111111)
	h := hash(id.Hash())
	key := engine.NewObjectIndex(id.Hash())
	i := indexFor(h, c.capacity)
	es := c.table[i]
	if es == nil {

		es = &entrys{
			size: 0,
		}
		var counter uint32 = (es.size + 1) << 16
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		counter += uint32(r.Int31n(1<<16 - 1))
		key.SetRand(counter)

		n := node.New(id)
		n.Index = key

		ne := &entry{
			hashCode: h,
			key:      key.LocalIndex(),
			value:    n,
			next:     nil,
		}
		es.head = ne
		c.table[i] = es

		c.size++
		if c.size >= c.threshold {
			c.resize(2 * c.capacity)
		}
		return n, true
	} else {
		for e := es.head; e != nil; e = e.next {
			if e.value.Id.Equals(id) {
				return e.value, false
			}
		}
		var counter uint32 = (es.size + 1) << 16
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		counter += uint32(r.Int31n(1<<16 - 1))
		key.SetRand(counter)

		n := node.New(id)
		n.Index = key

		ne := &entry{
			hashCode: h,
			key:      key.LocalIndex(),
			value:    n,
			next:     es.head,
		}
		es.head = ne
		es.size++

		c.size++
		if c.size >= c.threshold {
			c.resize(2 * c.capacity)
		}
		return n, true
	}

}

// Rehashes the contents of this map uint32o a new array with a
// larger capacity.  This method is called automatically when the
// number of keys in this map reaches its threshold.
//
// If current capacity is MAXIMUM_CAPACITY, this method does not
// resize the map, but sets threshold to Integer.MAX_VALUE.
// This has the effect of preventing future calls.
//
// newCapacity is the new capacity, MUST be a power of two;
// must be greater than current capacity unless current
// capacity is MAXIMUM_CAPACITY (in which case value
// is irrelevant).
func (c *Collection) resize(newCapacity uint32) {
	oldCapacity := c.capacity
	if oldCapacity == maximumCapacity {
		c.threshold = maximumCapacity
		return
	}

	newTable := make([]*entrys, newCapacity)
	c.transfer(newTable)
	c.table = newTable
	c.capacity = newCapacity
	c.threshold = uint32(float64(newCapacity) * c.loadFactor)
}

// Transfers all entries from current table to newTable.
func (c *Collection) transfer(newTable []*entrys) {
	src := c.table
	newCapacity := uint32(len(newTable))
	for j := 0; j < len(src); j++ {
		es := src[j]
		if es != nil {
			src[j] = nil
			e := es.head
			for e != nil {
				next := e.next
				i := indexFor(e.hashCode, newCapacity)
				es := newTable[i]
				if es == nil {
					e.next = nil
					nes := &entrys{
						head: e,
						size: 1,
					}
					newTable[i] = nes
				} else {
					e.next = es.head
					es.head = e
					es.size++
				}
				e = next
			}
		}
	}
}

// Removes the mapping for the specified key from this map if present.
//
// return if there was a mapping for key and
// the previous value associated with key, or
func (c Collection) Remove(key *engine.ObjectIndex) (*node.Node, bool) {
	e, found := c.removeEntryForKey(key)
	if found {
		return e.value, true
	}
	return nil, false
}

// Removes and returns the entry associated with the specified key
// in the Collection.
//
// Returns if there was a mapping for key and the associated entry
func (c Collection) removeEntryForKey(key *engine.ObjectIndex) (*entry, bool) {
	//hashCode := 0
	//if key != nil {
	//hashCode = hash(key.HashCode())
	//}
	hash := hash(key.GetId())
	i := indexFor(hash, c.capacity)
	es := c.table[i]
	pre := es.head
	for e := es.head; e != nil; e = e.next {
		if e.hashCode == hash && e.key == key.LocalIndex() {
			pre.next = e.next
			es.size--
			c.size--
		} //if
		pre = e
	}
	return nil, false
}

// Removes all of the mappings from this map.
// The map will be empty after this call returns.
func (c Collection) Clear() {
	tab := c.table
	var i uint32
	for i = 0; i < c.capacity; i++ {
		tab[i] = nil
	}
	c.size = 0
}

//Iter iterator all the items
func (c Collection) Iter() <-chan *node.Node {
	v := make(chan *node.Node)
	go func(v chan<- *node.Node) {
		var i uint32
		for i = 0; i < c.capacity; i++ {
			es := c.table[i]
			if es != nil {
				for e := es.head; e != nil; e = e.next {
					v <- e.value
				}
			}
		}
		close(v)
	}(v)
	return v
}
