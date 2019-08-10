package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type QueueRecord struct {
	ID   int
	Link string
}

func InsertQueueRecords(domainId int, urls []string) error {
	if 0 == len(urls) {
		return errors.New("InsertQueueRecords() no urls provided")
	}

	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO queue_record (domain_id, url) VALUES ")
	for _, url := range urls {
		buffer.WriteString("(")
		buffer.WriteString(strconv.Itoa(domainId))
		buffer.WriteString(",")
		buffer.WriteString(quoteString(url))
		buffer.WriteString("),")
	}
	buffer.Truncate(buffer.Len() - 1)

	tab, err := pool.Exec(context.Background(), buffer.String())
	if err != nil {
		return err
	}
	if tab.RowsAffected() != int64(len(urls)) {
		return errors.New(fmt.Sprintf("expected %d queue records to be inserted, %d were inserted", int64(len(urls)), tab.RowsAffected()))
	}

	return nil
}
