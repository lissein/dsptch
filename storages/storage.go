package storages

// Storage is the interface for storage backends (RDS, REDIS, DynamoDB, ...)
type Storage interface {
	save(id int)
}
