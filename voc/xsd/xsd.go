// Package xsd contains constants of the W3C XML Schema Definition Language https://www.w3.org/TR/xmlschema11-1/
package xsd

import "github.com/cayleygraph/quad/voc"

func init() {
	voc.RegisterPrefix(Prefix, NS)
}

const (
	NS     = "http://www.w3.org/2001/XMLSchema#"
	Prefix = "xsd"
)

// Base types
const (
	// Boolean represents the values of two-valued logic.
	Boolean = Prefix + `boolean`
	// String represents character strings
	String = Prefix + `string`
	// Double datatype is patterned after the IEEE double-precision 64-bit floating point datatype [IEEE 754-2008]. Each floating point datatype has a value space that is a subset of the rational numbers.  Floating point numbers are often used to approximate arbitrary real numbers.
	Double = Prefix + `double`
	// DateTime represents instants of time, optionally marked with a particular time zone offset.  Values representing the same instant but having different time zone offsets are equal but not identical.
	DateTime = Prefix + `dateTime`
)

// Extra numeric types
const (
	// Integer is derived from decimal by fixing the value of fractionDigits to be 0 and disallowing the trailing decimal point. This results in the standard mathematical concept of the integer numbers.
	Integer = Prefix + `integer`
	// Long is derived from integer by setting the value of maxInclusive to be 9223372036854775807 and minInclusive to be -9223372036854775808. The base type of long is integer.
	Long = Prefix + `long`
	// Int is derived from long by setting the value of maxInclusive to be 2147483647 and minInclusive to be -2147483648.
	Int = Prefix + `int`
	// Float datatype is patterned after the IEEE single-precision 32-bit floating point datatype
	Float = Prefix + `float`
)
