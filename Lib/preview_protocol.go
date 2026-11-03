package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const previewMagic = "OLCP0001"
const previewMaxPayload = 4 * 1024 * 1024

type previewRecord struct {
	kind                   byte
	key                    bool
	generation             uint32
	pts, dts               int64
	numerator, denominator uint32
	data                   []byte
}
type previewConfig struct {
	width, height, fpsNum, fpsDen uint32
	headers                       []byte
}

func readPreviewRecord(r io.Reader) (previewRecord, error) {
	var p previewRecord
	var length [4]byte
	if _, err := io.ReadFull(r, length[:]); err != nil {
		return p, err
	}
	n := binary.BigEndian.Uint32(length[:])
	if n < 36 || n > previewMaxPayload+36 {
		return p, fmt.Errorf("invalid preview record size")
	}
	var h [36]byte
	if _, err := io.ReadFull(r, h[:]); err != nil {
		return p, err
	}
	p.kind, p.key = h[0], h[1] == 1
	p.generation = binary.BigEndian.Uint32(h[4:8])
	p.pts = int64(binary.BigEndian.Uint64(h[8:16]))
	p.dts = int64(binary.BigEndian.Uint64(h[16:24]))
	p.numerator = binary.BigEndian.Uint32(h[24:28])
	p.denominator = binary.BigEndian.Uint32(h[28:32])
	size := binary.BigEndian.Uint32(h[32:36])
	if (p.kind != 2 && p.kind != 3) || h[1] > 1 || h[2] != 0 || h[3] != 0 ||
		p.generation == 0 || p.numerator == 0 || p.denominator == 0 || size != n-36 || size == 0 ||
		(p.kind == 2 && (p.key || size > 65536+16)) {
		return p, fmt.Errorf("invalid preview header")
	}
	p.data = make([]byte, size)
	_, err := io.ReadFull(r, p.data)
	return p, err
}

func parsePreviewConfig(p previewRecord) (previewConfig, error) {
	var c previewConfig
	if p.kind != 2 || len(p.data) < 17 {
		return c, fmt.Errorf("missing preview config")
	}
	c.width = binary.BigEndian.Uint32(p.data[0:4])
	c.height = binary.BigEndian.Uint32(p.data[4:8])
	c.fpsNum = binary.BigEndian.Uint32(p.data[8:12])
	c.fpsDen = binary.BigEndian.Uint32(p.data[12:16])
	if c.width < 2 || c.height < 2 || c.width > 1280 || c.height > 720 || c.width%2 != 0 || c.height%2 != 0 ||
		c.fpsNum == 0 || c.fpsDen == 0 || uint64(c.fpsNum) > uint64(c.fpsDen)*30 {
		return c, fmt.Errorf("unsupported preview format")
	}
	c.headers = p.data[16:]
	sps, pps := false, false
	for _, nal := range annexBNALs(c.headers) {
		switch nal[0] & 31 {
		case 7:
			sps = true
		case 8:
			pps = true
		default:
			return c, fmt.Errorf("invalid H264 initialization NAL")
		}
	}
	if !sps || !pps {
		return c, fmt.Errorf("missing Annex-B SPS/PPS")
	}
	return c, nil
}

// Return NALs only for an Annex-B byte stream, never guess AVCC lengths.
func annexBNALs(data []byte) [][]byte {
	var result [][]byte
	start := -1
	for i := 0; i+2 < len(data); {
		size := 0
		if data[i] == 0 && data[i+1] == 0 {
			if data[i+2] == 1 {
				size = 3
			} else if i+3 < len(data) && data[i+2] == 0 && data[i+3] == 1 {
				size = 4
			}
		}
		if size == 0 {
			i++
			continue
		}
		if start == -1 && !bytes.Equal(data[:i], make([]byte, i)) {
			return nil
		}
		if start >= 0 && i > start {
			result = append(result, data[start:i])
		}
		start = i + size
		i += size
	}
	if start >= 0 && start < len(data) {
		result = append(result, data[start:])
	}
	return result
}

func validPreviewFrame(p previewRecord) bool {
	nals := annexBNALs(p.data)
	if len(nals) == 0 {
		return false
	}
	idr, slice := false, false
	for _, n := range nals {
		if n[0]&0x80 != 0 {
			return false
		}
		kind := n[0] & 31
		if kind == 5 {
			idr = true
		}
		if kind == 1 || kind == 5 {
			slice = true
		}
	}
	return slice && idr == p.key
}
