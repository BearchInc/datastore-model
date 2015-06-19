package db

type SoftDeletableModel struct {
	Model
	Deleted bool `json:"-"`
}

// Implements SoftDeletableEntity
func (this *SoftDeletableModel) IsDeleted() bool {
	return this.Deleted
}