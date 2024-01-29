package adapter

// seqGT32 true if a > b
func seqGT32(a, b uint32) bool {
	return int32(a-b) > 0
}

// setGT16 true if a > b
func setGT16(a, b uint16) bool {
	return int16(a-b) > 0
}
