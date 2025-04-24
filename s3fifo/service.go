package s3fifo

//
//import "s3fifo/structures"
//
//type key struct {
//	hash  uint64
//	count uint16
//}
//type Value[V any] struct {
//	v       V
//	count   int
//	isGhost bool
//}
//type Cache[V any] struct {
//	main  structures.SimpleQueue[structures.Node[V]]
//	small structures.SimpleQueue[structures.Node[V]]
//
//	ghost structures.SimpleQueue[uint64]
//
//	data map[uint64]*Value[V]
//}
//
//func (m *Cache[V]) Insert(k uint64, v V) {
//	value, ok := m.data[k]
//	if ok {
//		if value.isGhost {
//			if m.main.IsFull() {
//				evictedHash := m.main.Pop()
//				evicted, okEvicted := m.data[evictedHash]
//				if !okEvicted {
//					panic("")
//				}
//
//				if evicted.count != 0 {
//					//evicted.count--
//					//m.Insert(evictedHash, evicted.v)
//					// TODO: reinsert.
//				}
//			}
//			m.main.Put(k)
//			value.isGhost = false
//			value.v = v
//		} else {
//			value.count++
//		}
//	} else {
//		if m.small.IsFull() {
//			evictedHash := m.small.Pop()
//			evicted, okEvicted := m.data[evictedHash]
//			if !okEvicted {
//				panic("")
//			}
//
//			if evicted.count == 0 {
//				m.ghost.Put(evictedHash)
//			} else {
//				evicted.count--
//				if m.main.IsFull() {
//					evictedHashmain := m.main.Pop()
//					evictedMain, okEvictedMain := m.data[evictedHashmain]
//					if !okEvictedMain {
//						panic("")
//					}
//
//					if evictedMain.count != 0 {
//						//evicted.count--
//						//m.Insert(evictedHash, evicted.v)
//						// TODO: reinsert.
//					}
//				}
//				m.main.Put(evictedHash)
//			}
//		}
//
//		m.small.Put(k)
//		m.data[k] = &Value[V]{
//			v:       v,
//			count:   0,
//			isGhost: false,
//		}
//	}
//
//	return
//
//}
