package pager

import (
	"io"
	"os"
)

// Pager uses a basic caching mechanism to store pages in memory once
// they've been accessed. It makes no attempt at any other optimizations.
type Pager struct {
	File   *os.File
	Cache  map[uint64][]byte
	Config PagerConfig
}

type PagerConfig struct {
	PageSize      uint64
	ReservedSpace uint64
}

func NewPager(filepath string, config PagerConfig) (*Pager, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &Pager{
		File:   file,
		Config: config,
		Cache:  make(map[uint64][]byte),
	}, nil
}

func (p *Pager) Close() error {
	return p.File.Close()
}

func (p *Pager) PageSize() uint64 {
	return p.Config.PageSize
}

func (p *Pager) ReservedSpace() uint64 {
	return p.Config.ReservedSpace
}

func (p *Pager) GetPage(pageNum uint64) ([]byte, error) {
	// If the page isn't already in the cache, we should read it directly from
	// the file
	if _, ok := p.Cache[pageNum]; !ok {
		// Calculate the byte number at which this page starts
		pageStart := int64((pageNum - 1) * uint64(p.PageSize()))

		// Seek to the beginning of the page
		p.File.Seek(pageStart, io.SeekStart)

		// Read the entire page into a byte slice
		pageBytes := make([]byte, p.PageSize())
		if _, err := p.File.Read(pageBytes); err != nil {
			return nil, err
		}

		// Insert the page into cache
		p.Cache[pageNum] = pageBytes
	}

	return p.Cache[pageNum], nil
}
