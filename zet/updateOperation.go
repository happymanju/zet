package zet

type UpdateOperation struct {
	ToAdd    []string
	ToRemove []string
}

func newUpdateOperation() UpdateOperation {
	u := UpdateOperation{
		ToAdd:    []string{},
		ToRemove: []string{},
	}

	return u
}
