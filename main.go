package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

func parsePoint(s string) (Point, error) {
	// Split the line on the comma...
	sns := strings.Split(s, ",")
	if len(sns) != 2 {
		return Point{}, fmt.Errorf("invalid line format %q", s)
	}

	// Parse the x and y values as integers...
	x, err := strconv.Atoi(sns[0])
	if err != nil {
		return Point{}, fmt.Errorf("failed to parse x value as int %q", sns[0])
	}
	y, err := strconv.Atoi(sns[1])
	if err != nil {
		return Point{}, fmt.Errorf("failed to parse y value as int %q", sns[1])
	}

	// Return the parsed point...
	return Point{X: x, Y: y}, nil
}

func parseFile(f *os.File) ([]Point, error) {
	// Create a scanner for the file...
	buf := bufio.NewScanner(f)

	// Create a slice of points to return...
	var points []Point

	// Iterate over the lines...
	for buf.Scan() {
		// Parse the point...
		p, err := parsePoint(buf.Text())
		if err != nil {
			return nil, err
		}

		// Append the point to the slice...
		points = append(points, p)
	}

	// Return the slice of points...
	return points, nil
}

func dist(a, b Point) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func closest(p Point, cs []Point) int {
	besti, bestn := -1, 0
	for i, c := range cs {
		n := dist(p, c)
		if besti == -1 || n < bestn {
			besti, bestn = i, n
		}
	}
	return besti
}

func sortCentroids(ps []Point) {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].X < ps[j].X {
			return true
		}
		if ps[i].X > ps[j].X {
			return false
		}
		return ps[i].Y < ps[j].Y
	})
}

func compareCentroids(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i, p := range a {
		if p != b[i] {
			return false
		}
	}
	return true
}

func main() {
	// Load the number of clusters...
	pk := flag.Int("k", 3, "number of clusters")
	flag.Parse()
	k := *pk

	// Load the points from stdin...
	points, err := parseFile(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Pick the first k points as the
	// initial centroids...
	centroids := make([]Point, k)
	for i := 0; i < k; i++ {
		centroids[i] = points[i]
	}

	// Sort the centroids...
	sortCentroids(centroids)

	// Start looping...
	for {
		// Find the closest centroid for each point...
		closestCentroids := make([]int, len(points))
		for i, p := range points {
			closestCentroids[i] = closest(p, centroids)
		}

		// Create a slice to store the new centroids...
		newCentroids := make([]Point, k)
		counts := make([]int, k)
		for i, p := range points {
			// Get that centroid...
			ci := closestCentroids[i]

			// Add the point to the centroid...
			newCentroids[ci].X += p.X
			newCentroids[ci].Y += p.Y

			// Increment the count...
			counts[ci]++
		}

		// Average the centroids...
		for i, n := range counts {
			newCentroids[i].X /= n
			newCentroids[i].Y /= n
		}

		// Sort the new centroids...
		sortCentroids(newCentroids)

		// Check if the centroids have changed...
		if compareCentroids(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}

	// Print the centroids...
	for _, c := range centroids {
		fmt.Printf("%d,%d\n", c.X, c.Y)
	}
}
