package pagemanager

import (
	"encoding/binary"
	"hash/crc64"
)

const (
	PageSize       = 4096
	PageHeaderSize = 64
	PageDataSize   = PageSize - PageHeaderSize
)

// PageID represents a unique page identifier
type PageID uint64

// Page represents a single page in the database
type Page struct {
	ID       PageID
	Data     [PageDataSize]byte
	Dirty    bool
	Pinned   bool
	checksum uint64
}

// NewPage creates a new page with the given ID
func NewPage(id PageID) *Page {
	return &Page{
		ID:     id,
		Dirty:  false,
		Pinned: false,
	}
}

// ComputeChecksum computes the CRC64 checksum of the page data
func (p *Page) ComputeChecksum() uint64 {
	table := crc64.MakeTable(crc64.ISO)
	return crc64.Checksum(p.Data[:], table)
}

// Validate checks if the page checksum is valid
func (p *Page) Validate() bool {
	return p.checksum == p.ComputeChecksum()
}

// Marshal serializes the page to bytes
func (p *Page) Marshal() []byte {
	buf := make([]byte, PageSize)

	// Header
	binary.LittleEndian.PutUint64(buf[0:8], uint64(p.ID))
	binary.LittleEndian.PutUint64(buf[8:16], p.ComputeChecksum())

	// Data
	copy(buf[PageHeaderSize:], p.Data[:])

	return buf
}

// Unmarshal deserializes bytes into a page
func (p *Page) Unmarshal(data []byte) error {
	// TODO: Implement deserialization
	// - Read header fields
	// - Copy data
	// - Validate checksum
	return nil
}
