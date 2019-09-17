package order

func addNewElevatorsToDatabase(database [][][]bool, numberOfNewElevators int, numberOfFloors int) [][][]bool {
	for i := 0; i < numberOfNewElevators; i++ {
		database = append(database, CreateEmptyQueue(numberOfFloors))
	}
	return database
}