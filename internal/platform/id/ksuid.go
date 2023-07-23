package id

import (
	"fmt"
	"github.com/segmentio/ksuid"
)

type KsuidGenerator struct{}

func NewKsuidGenerator() *KsuidGenerator {
	return &KsuidGenerator{}
}

func (k *KsuidGenerator) MustGenerate() string {
	id, err := ksuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (k *KsuidGenerator) Generate() (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrorGeneratingID, err)
	}
	return id.String(), nil
}
