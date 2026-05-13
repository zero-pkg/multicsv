package multicsv

import (
	"bytes"
	"encoding/csv"
	"io"
)

const detectDelimiterSampleSize = 32 * 1024

var delimiterCandidates = []rune{',', ';', '\t', '|'}

type delimiterScore struct {
	consistentRecords int
	fieldCount        int
	totalFields       int
	errors            int
}

func detectDelimiter(r io.ReadSeeker) (rune, error) {
	sample := make([]byte, detectDelimiterSampleSize)
	n, err := r.Read(sample)
	seekErr := resetReader(r)

	if err != nil && err != io.EOF {
		return 0, err
	}

	if seekErr != nil {
		return 0, seekErr
	}

	return detectDelimiterFromSample(sample[:n]), nil
}

func resetReader(r io.Seeker) error {
	_, err := r.Seek(0, io.SeekStart)
	return err
}

func detectDelimiterFromSample(sample []byte) rune {
	bestDelimiter := ','
	bestScore := delimiterScore{}

	for _, delimiter := range delimiterCandidates {
		score := scoreDelimiter(sample, delimiter)
		if isBetterDelimiterScore(score, bestScore) {
			bestDelimiter = delimiter
			bestScore = score
		}
	}

	return bestDelimiter
}

func scoreDelimiter(sample []byte, delimiter rune) delimiterScore {
	reader := csv.NewReader(bytes.NewReader(sample))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	fieldCounts := make(map[int]int)
	score := delimiterScore{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			score.errors++
			break
		}

		fieldCount := len(record)
		if fieldCount <= 1 {
			continue
		}

		fieldCounts[fieldCount]++

		score.totalFields += fieldCount
		if fieldCounts[fieldCount] > score.consistentRecords {
			score.consistentRecords = fieldCounts[fieldCount]
			score.fieldCount = fieldCount
		}
	}

	return score
}

func isBetterDelimiterScore(score, best delimiterScore) bool {
	if score.consistentRecords != best.consistentRecords {
		return score.consistentRecords > best.consistentRecords
	}

	if score.fieldCount != best.fieldCount {
		return score.fieldCount > best.fieldCount
	}

	if score.totalFields != best.totalFields {
		return score.totalFields > best.totalFields
	}

	return score.errors < best.errors
}
