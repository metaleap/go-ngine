//	Loads assets from Collada 1.4.1 and 1.5 documents into the data structures provided by the go:ngine *assets* package.
//	Note, Collada 1.4.1 documents are insta-bootstrapped to version 1.5 in-memory.
package collimp

import (
	xsdt "github.com/metaleap/go-xsd/types"
)

func f64(d xsdt.ToXsdtDouble) float64 {
	return float64(d.ToXsdtDouble())
}
