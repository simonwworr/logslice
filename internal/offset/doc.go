// Package offset provides a lightweight byte-offset index for log files.
//
// An Index records (lineNumber, byteOffset) pairs at regular intervals so that
// a reader can seek close to a target line without scanning from the beginning
// of the file. This is especially valuable for large, multi-gigabyte logs where
// even a modest index density (e.g. every 1 000 lines) reduces seek cost by
// several orders of magnitude.
//
// Typical usage:
//
//	idx := offset.New()
//	var pos int64
//	for lineNum, line := range lines {
//		if lineNum%1000 == 0 {
//			idx.Add(int64(lineNum), pos)
//		}
//		pos += int64(len(line))
//	}
//	idx.SaveToFile("mylog.idx")
//
// To seek:
//
//	idx, _ := offset.LoadFromFile("mylog.idx")
//	rec, ok := idx.Lookup(targetLine)
//	if ok {
//		file.Seek(rec.ByteOffset, io.SeekStart)
//	}
package offset
