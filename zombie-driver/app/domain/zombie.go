package domain

type Zombier interface {
	IsZombie(driverDistance float64) bool
}

type ZombieDetector struct {
	MinimumDistance float64
}

func (z ZombieDetector) IsZombie(driverDistance float64) bool {
	return driverDistance < z.MinimumDistance
}
