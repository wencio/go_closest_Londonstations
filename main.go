package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"

	"github.com/gocarina/gocsv"
)

type Station struct {
	Name      string  `csv:"Station"`
	Latitude  float64 `csv:"Latitude"`
	Longitude float64 `csv:"Longitude"`
}

type Position struct {
	Longitude float64
	Latitude  float64
}

type OutputStation struct {
	Name     string  `csv:"station_name"`
	Distance float64 `csv:"distance"`
}

func main() {

	bytes, err := ioutil.ReadFile("London.csv")

	if err != nil {
		panic(err)

	}

	var stations []Station

	_ = gocsv.UnmarshalBytes(bytes, &stations)

	target := Position{Latitude: 51.479495, Longitude: -0.000500}

	outputStations := make([]OutputStation, len(stations))

	for i, station := range stations {

		outputStations[i] = OutputStation{

			Name:     station.Name,
			Distance: distanceInKmBetweenEarthCoordinates(target.Latitude, target.Longitude, station.Latitude, station.Longitude),
		}

	}

	sort.SliceStable(outputStations, func(i, j int) bool {
		return outputStations[i].Distance < outputStations[j].Distance

	})

	log.Print(outputStations[0])

	file, _ := os.Create("closest_stations.csv")

	gocsv.Marshal(&outputStations, file)

}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func distanceInKmBetweenEarthCoordinates(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	var earthRadiusKm = 6371.0

	var dLat = degreesToRadians(lat2 - lat1)
	var dLon = degreesToRadians(lon2 - lon1)

	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}
