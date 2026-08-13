package mongo

func fromMongoErr(err error) error {
	if err == nil {
		return nil
	}
	// TODO: Complete impl
	return err
}