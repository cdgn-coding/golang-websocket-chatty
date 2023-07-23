package id

import "errors"

var ErrorGeneratingID = errors.New("error generating id")

type Generator interface {
	MustGenerate() string
	Generate() (string, error)
}
