package store_redis

type Redis interface {
	Inspection() InspectionRepository
}
