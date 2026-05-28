package adapters

import (
	"context"
	"sync"
)

type SpreadsheetServiceStub struct {
	lock sync.Mutex

	Rows []string
}

func (s *SpreadsheetServiceStub) AppendRow(ctx context.Context, sheetName string, row []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Rows = append(s.Rows, row...)

	return nil
}
