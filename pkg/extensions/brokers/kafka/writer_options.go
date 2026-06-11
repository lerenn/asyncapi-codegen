package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// WriterOption is a function that configures a kafka.Writer.
type WriterOption func(*kafka.Writer)

// WithMaxAttempts limit on how many attempts will be made to deliver a message.
//
// The default is to try at most 10 times.
func WithMaxAttempts(maxAttempts int) WriterOption {
	return func(w *kafka.Writer) {
		w.MaxAttempts = maxAttempts
	}
}

// WithWriteBackoffMin optionally sets the smallest amount of time the writer waits before
// it attempts to write a batch of messages.
//
// Default: 100ms.
func WithWriteBackoffMin(v time.Duration) WriterOption {
	return func(w *kafka.Writer) {
		w.WriteBackoffMin = v
	}
}

// WithWriteBackoffMax optionally sets the maximum amount of time the writer waits before
// it attempts to write a batch of messages.
//
// Default: 1s.
func WithWriteBackoffMax(v time.Duration) WriterOption {
	return func(w *kafka.Writer) {
		w.WriteBackoffMax = v
	}
}

// WithBatchSize limits on how many messages will be buffered before being sent to a
// partition.
//
// The default is to use a target batch size of 100 messages.
func WithBatchSize(v int) WriterOption {
	return func(w *kafka.Writer) {
		w.BatchSize = v
	}
}

// WithBatchBytes limits the maximum size of a request in bytes before being sent to
// a partition.
//
// The default is to use a kafka default value of 1048576.
func WithBatchBytes(v int64) WriterOption {
	return func(w *kafka.Writer) {
		w.BatchBytes = v
	}
}

// WithBatchTimeout is a timeout on how often incomplete message batches will be flushed to
// kafka.
//
// The default is to flush at least every second.
func WithBatchTimeout(v time.Duration) WriterOption {
	return func(w *kafka.Writer) {
		w.BatchTimeout = v
	}
}

// WithReadTimeout is a timeout for read operations performed by the Writer.
//
// Defaults to 10 seconds.
func WithReadTimeout(v time.Duration) WriterOption {
	return func(w *kafka.Writer) {
		w.ReadTimeout = v
	}
}

// WithWriteTimeout is a timeout for write operation performed by the Writer.
//
// Defaults to 10 seconds.
func WithWriteTimeout(v time.Duration) WriterOption {
	return func(w *kafka.Writer) {
		w.WriteTimeout = v
	}
}

// WithAsync Setting this flag to true causes the WriteMessages method to never block.
// It also means that errors are ignored since the caller will not receive
// the returned value. Use this only if you don't care about guarantees of
// whether the messages were written to kafka.
//
// Defaults to false.
func WithAsync(v bool) WriterOption {
	return func(w *kafka.Writer) {
		w.Async = v
	}
}
