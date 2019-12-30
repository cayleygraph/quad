// Package owl contains constants of the Web Ontology Language (OWL)
package owl

import "github.com/cayleygraph/quad/voc"

func init() {
	voc.RegisterPrefix(Prefix, NS)
}

const (
	NS     = `http://www.w3.org/2002/07/owl#`
	Prefix = `owl:`
)

const (
	UnionOf          = Prefix + "unionOf"
	Restriction      = Prefix + "Restriction"
	OnProperty       = Prefix + "onProperty"
	Cardinality      = Prefix + "cardinality"
	MaxCardinality   = Prefix + "maxCardinality"
	Thing            = Prefix + "Thing"
	Class            = Prefix + "Class"
	DatatypeProperty = Prefix + "DatatypeProperty"
	ObjectProperty   = Prefix + "ObjectProperty"
)
