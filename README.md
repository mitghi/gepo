# gepo

[<img src="https://img.shields.io/badge/godocs-gepo-mitghi"></img>](https://pkg.go.dev/github.com/mitghi/gepo)
![GitHub](https://img.shields.io/github/license/mitghi/gepo)

Gepo is simple implementation of NNS (Proximity search) algorithm. It provides minimal, simple and easy to use api 
to insert and search nearest points on a plane which is divided into tiles.

The library has no dependencies other than standard library.

```
$ go get -u github.com/mitghi/gepo
```

following is an example of using the api:

```go
package main

import (
	"fmt"
	"github.com/mitghi/gepo"
)

func main() {
	// initialize the plane
	// with resolution of 20 Kilometer
	plane := gepo.New(gepo.Km(20))

	berlinMitte := gepo.NewPoint("Berlin Mitte", 52.519294, 13.405868)
	berlinPanaromaStr := gepo.NewPoint("Berlin Panaroma str", 52.520783, 13.409578)
	berlinMall := gepo.NewPoint("Berlin Mall", 52.510415, 13.381302)
	berlinPostdamerPlatz := gepo.NewPoint("Berlin Postdamer Platz", 52.509678, 13.375129)

	// add points to the plane
	// the points are inserted
	// based on tile division
	plane.AddPoints([]*gepo.Point{
		berlinMitte,
		berlinPanaromaStr,
		berlinMall,
		berlinPostdamerPlatz}...)

	// this closure is for extra control
	// over what gets included in the final
	// results; that is the user can decide
	// whether certain points are unfit for
	// inclusion
	justAccept := func(_ *gepo.Point) bool { return true }

	// perform serach

	// returns an slice containing (berlinMitte, berlinPanaromaStr)
	neighbours := plane.Nearest(berlinMitte, 2000, gepo.Km(1), justAccept)
	fmt.Println("from berlinMitte:", neighbours)
	// ...

	// return (berilnPanaromaStr, berlinMitte, berlinMall)
	neighbours = plane.Nearest(berlinPanaromaStr, 2000, gepo.Km(3.4), justAccept)
	fmt.Println("from berlinPanaromaStr:", neighbours)
	// ...

	// return (berlinMall, berilnPostdamerPlatz)
	neighbours = plane.Nearest(berlinMall, 2000, gepo.Km(1), justAccept)
	fmt.Println("from berlinMall:", neighbours)
	// ...
}
```
