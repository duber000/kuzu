package memoryleak

// BUGGY: MVCC with memory leak
type Version struct {
	data []byte
	prev *Version  // LEAK: never cleaned up
}

type MVCCStore struct {
	versions map[string]*Version
}

func (s *MVCCStore) Write(key string, data []byte) {
	// LEAK: Old versions accumulate!
	newVer := &Version{
		data: data,
		prev: s.versions[key],
	}
	s.versions[key] = newVer
}

// BUGGY: Buffer Pool with memory leak
type BufferPool struct {
	frames    []*Frame
	pageTable map[int]int
}

type Frame struct {
	data []byte
}

func (bp *BufferPool) Evict(frameID int) {
	// LEAK: Frame not removed from frames array!
	delete(bp.pageTable, frameID)
}

// BUGGY: Query result leak
type ResultSet struct {
	rows []Row
	pos  int
}

type Row map[string]interface{}

func (rs *ResultSet) Next() bool {
	rs.pos++
	return rs.pos < len(rs.rows)
}

func (rs *ResultSet) Close() {
	// Should free resources
	rs.rows = nil
}

type Database struct{}

func (db *Database) ExecuteQuery(sql string) error {
	rs := &ResultSet{
		rows: make([]Row, 1000000),  // Large allocation
	}

	for rs.Next() {
		// Process
	}

	// LEAK: rs.Close() never called!
	return nil
}

// TODO: Fix these memory leaks!
// Hints:
// 1. Implement GC for old versions
// 2. Clear frame references
// 3. Always close result sets (use defer)
// 4. Profile with pprof to verify fixes
