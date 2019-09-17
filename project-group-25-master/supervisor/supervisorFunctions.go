package supervisor

//appendElevators append new elevators to the alive map
func appendElevators(listOfElev map[int]bool, numberOfNewElevators int) map[int]bool {
	for i := 0; i < numberOfNewElevators; i++ {
		listOfElev[i] = false
	}
	return listOfElev
}

//createAlive creates the alive map for the first active elevator
func createAlive(elevatorID int) map[int]bool {
	alive := make(map[int]bool)
	for i := 0; i < elevatorID+1; i++ {
		alive[i] = false
	}
	return alive
}
