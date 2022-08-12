package bulker

import "github.com/jitsucom/bulker/types"

type BulkMode int

const (
	//AutoCommit - bulker stream immediately commits each consumed object to the database
	//Useful when working with live stream of objects
	AutoCommit BulkMode = iota

	//Transactional - bulker stream commits all consumed object on Complete call
	//Useful when working with large number of objects or to optimize performance or costs using batch processing
	//Any error with just one object will fail the whole transaction and no objects will be written to the database
	Transactional

	//ReplacePartition stream replaces all rows associated with the chosen partition column value in a single transaction (where applicable).
	//It is useful when it is required to reprocess all objects associates with specific partition id.
	//E.g. for processing and reprocessing one day of reporting data
	//If data of your stream may be reprocessed in some point in time it is recommended to always use ReplacePartition mode for that stream
	//
	//ReplacePartition implies Transactional, meaning that the new data will be available only after BulkerStream.complete() call
	ReplacePartition

	//ReplaceTable - atomically replaces target table with a new one filled with the object injected to current stream.
	//To sync entire collection of object at once without leaving target table in unfinished state
	//Useful when collection contains finite number of object, and when it is required that target table always represent complete state for some point of time.
	//
	//ReplaceTable implies Transactional, meaning that the new data will be available only after BulkerStream.complete() call
	ReplaceTable
)

// Bulker interface allows streaming object to data warehouses using different modes.
// See BulkMode for more details.
type Bulker interface {
	//CreateStream create a BulkerStream instance that will store objects to the target table in a data warehouse.
	//bulker BulkerStream creates a new table with provided tableName if it does not exist.
	//Table schema is based on flattened object structure but may be overridden by providing WithTable option.
	//bulker BulkerStream will add new column to a table on the fly if new properties appear in object and table schema is not overridden.
	CreateStream(id, tableName string, mode BulkMode, streamOptions ...StreamOption) (BulkerStream, error)
}

// TODO: Commit() method that commits transaction and start new one ??
type BulkerStream interface {
	//Consume - put object to the stream. If stream is in AutoCommit mode it will be immediately committed to the database.
	//Otherwise, it will be buffered and committed on Complete call.
	Consume(object types.Object) error
	//Abort - abort stream and rollback all uncommitted objects. For stream in AutoCommit mode does nothing.
	//Returns stream statistics. BulkerStream cannot be used after Abort call.
	Abort() (State, error)
	//Complete - commit all uncommitted objects to the database. For stream in AutoCommit mode does nothing.
	//Returns stream statistics. BulkerStream cannot be used after Complete call.
	Complete() (State, error)
}

type Status string

const (
	//Completed - stream was completed successfully
	Completed Status = "Completed"
	//Aborted - stream was aborted by user
	Aborted = "Aborted"
	//Failed - failed to complete stream
	Failed = "Failed"
	//Active - stream is active
	Active = "Active"
)

// State is used as a Batch storing result
type State struct {
	Status         Status
	LastError      error
	ProcessedRows  int
	SuccessfulRows int
	RowsErrors     map[int]error
}
