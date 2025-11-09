package mvccreview

// REVIEW THIS CODE

type Timestamp uint64
type TxnID uint64

type Version struct {
	data    []byte
	beginTS Timestamp
	endTS   *Timestamp
	prev    *Version  // Issue 1: Should use weak.Pointer
}

type MVCCStore struct {
	versions map[string]*Version
}

// Issue 2: Wrong visibility check
func (s *MVCCStore) isVisible(v *Version, snapshot Timestamp) bool {
	// WRONG: Should check v.endTS > snapshot, not >=
	return v.beginTS <= snapshot && (v.endTS == nil || *v.endTS >= snapshot)
}

// Issue 3: No write conflict detection
func (s *MVCCStore) Commit(txn TxnID, writes map[string][]byte) error {
	// WRONG: No conflict detection!
	for key, data := range writes {
		newVer := &Version{
			data:    data,
			beginTS: Timestamp(txn),
			prev:    s.versions[key],
		}
		s.versions[key] = newVer
	}
	return nil
}

// Issue 4: GC removes live versions
func (s *MVCCStore) GC(minSnapshot Timestamp) {
	for key, ver := range s.versions {
		// WRONG: May delete visible versions!
		if ver.beginTS < minSnapshot {
			delete(s.versions, key)
		}
	}
}

// Issue 5: Race condition in read
func (s *MVCCStore) Read(key string, snapshot Timestamp) []byte {
	// WRONG: No lock protection
	ver := s.versions[key]
	for ver != nil {
		if s.isVisible(ver, snapshot) {
			return ver.data
		}
		ver = ver.prev
	}
	return nil
}

// TODO: Find all MVCC bugs
