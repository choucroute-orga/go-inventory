package api

import (
	"fmt"
	"inventory/db"
	"inventory/messages"
)

// Conversion factors to base units
const (
	// Volume conversions (to ml)
	tspToMl  = 4.929
	tbspToMl = 14.787
	csToMl   = 250.0
	lToMl    = 1000.0

	// Weight conversions (to g)
	kgToG = 1000.0
)

// ConversionResult holds the converted quantity and unit
type ConversionResult struct {
	Quantity float64
	Unit     db.UnitType
}

// ConvertToBaseUnitFromRequest converts a quantity from a UnitRequest to its base UnitType
func ConvertToBaseUnitFromRequest(quantity float64, unit messages.UnitRequest) (ConversionResult, error) {
	switch unit {
	// Volume conversions
	case messages.UnitTsp:
		return ConversionResult{
			Quantity: quantity * tspToMl,
			Unit:     db.UnitMl,
		}, nil
	case messages.UnitTbsp:
		return ConversionResult{
			Quantity: quantity * tbspToMl,
			Unit:     db.UnitMl,
		}, nil
	case messages.UnitCs:
		return ConversionResult{
			Quantity: quantity * csToMl,
			Unit:     db.UnitMl,
		}, nil
	case messages.UnitL:
		return ConversionResult{
			Quantity: quantity * lToMl,
			Unit:     db.UnitMl,
		}, nil
	case messages.UnitMl:
		return ConversionResult{
			Quantity: quantity,
			Unit:     db.UnitMl,
		}, nil

	// Weight conversions
	case messages.UnitKg:
		return ConversionResult{
			Quantity: quantity * kgToG,
			Unit:     db.UnitG,
		}, nil
	case messages.UnitG:
		return ConversionResult{
			Quantity: quantity,
			Unit:     db.UnitG,
		}, nil

	// Item conversions
	case messages.UnitItem, messages.UnitItems:
		return ConversionResult{
			Quantity: quantity,
			Unit:     db.UnitItem,
		}, nil

	default:
		return ConversionResult{}, fmt.Errorf("unsupported unit request: %s", unit)
	}
}

func calculateConversionResult(quantity float64, ratio float64) ConversionResult {
	return ConversionResult{
		Quantity: quantity * ratio,
		Unit:     db.UnitG,
	}
}

// convertToBaseUnit converts a quantity to its base unit (g, ml or item)
func ConvertToBaseUnit(quantity float64, unit db.UnitType) (ConversionResult, error) {
	switch unit {
	// Volume conversions
	case "tsp":
		return ConversionResult{Quantity: quantity * tspToMl, Unit: db.UnitMl}, nil
	case "tbsp":
		return ConversionResult{Quantity: quantity * tbspToMl, Unit: db.UnitMl}, nil
	case "cs":
		return ConversionResult{Quantity: quantity * csToMl, Unit: db.UnitMl}, nil
	case "l":
		return ConversionResult{Quantity: quantity * lToMl, Unit: db.UnitMl}, nil
	case "ml":
		return ConversionResult{Quantity: quantity, Unit: db.UnitMl}, nil

	// Weight conversions
	case "kg":
		return ConversionResult{Quantity: quantity * kgToG, Unit: db.UnitG}, nil
	case "g":
		return ConversionResult{Quantity: quantity, Unit: db.UnitG}, nil

	// Item counts
	case "i", "is":
		return ConversionResult{Quantity: quantity, Unit: db.UnitItem}, nil

	default:
		return ConversionResult{0, ""}, fmt.Errorf("unsupported unit: %s", unit)
	}
}

// This function rounds the Base Unit if it's ml or g to a more readable unit
func roundBaseUnit(conversion ConversionResult) messages.ConversionMessageResult {
	quantity, unit := conversion.Quantity, conversion.Unit

	if unit == db.UnitMl && quantity >= lToMl {
		return messages.ConversionMessageResult{
			Quantity: quantity / lToMl,
			Unit:     messages.UnitMessageL,
		}
	}

	if unit == db.UnitG && quantity >= kgToG {
		return messages.ConversionMessageResult{
			Quantity: quantity / kgToG,
			Unit:     messages.UnitMessageKg,
		}
	}

	return messages.ConversionMessageResult{
		Quantity: quantity,
		Unit:     messages.UnitMessage(unit),
	}
}
