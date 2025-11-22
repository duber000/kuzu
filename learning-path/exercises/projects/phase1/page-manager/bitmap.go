package pagemanager

// Bitmap represents a bitmap for tracking free/allocated pages
type Bitmap struct {
	bits []byte
	size int
}

// NewBitmap creates a new bitmap with the given number of bits
func NewBitmap(size int) *Bitmap {
	numBytes := (size + 7) / 8
	return &Bitmap{
		bits: make([]byte, numBytes),
		size: size,
	}
}

// Set sets the bit at position n to 1
func (b *Bitmap) Set(n int) {
	if n < 0 || n >= b.size {
		return
	}
	b.bits[n/8] |= (1 << (n % 8))
}

// Clear sets the bit at position n to 0
func (b *Bitmap) Clear(n int) {
	if n < 0 || n >= b.size {
		return
	}
	b.bits[n/8] &= ^(1 << (n % 8))
}

// Test returns true if the bit at position n is 1
func (b *Bitmap) Test(n int) bool {
	if n < 0 || n >= b.size {
		return false
	}
	return (b.bits[n/8] & (1 << (n % 8))) != 0
}

// FindFirstZero finds the first 0 bit (free page)
func (b *Bitmap) FindFirstZero() int {
	// TODO: Implement efficient search
	// Optimization: scan bytes first, then bits
	for i := 0; i < b.size; i++ {
		if !b.Test(i) {
			return i
		}
	}
	return -1 // No free pages
}

// CountOnes returns the number of 1 bits (allocated pages)
func (b *Bitmap) CountOnes() int {
	// TODO: Implement using bit manipulation tricks
	count := 0
	for i := 0; i < b.size; i++ {
		if b.Test(i) {
			count++
		}
	}
	return count
}

// Resize grows or shrinks the bitmap
func (b *Bitmap) Resize(newSize int) {
	// TODO: Implement resize
	// - Allocate new byte slice
	// - Copy existing bits
	// - Update size
}
