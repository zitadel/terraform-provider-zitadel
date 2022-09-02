package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Stringified struct {
	Str string
}

func (s *Stringified) String() string {
	return s.Str
}

type Stringify interface {
	String() string
}

func SetToStringSlice(set *schema.Set) []string {
	slice := make([]string, 0)
	for _, secondFactor := range set.List() {
		slice = append(slice, secondFactor.(string))
	}
	return slice
}

func GetAddAndDelete(current []Stringify, desired []string) ([]string, []string) {
	addSlice := make([]string, 0)
	deleteSlice := make([]string, 0)

	for _, desiredItem := range desired {
		found := false
		for _, currentItem := range current {
			if desiredItem == currentItem.String() {
				found = true
			}
		}
		if !found {
			addSlice = append(addSlice, desiredItem)
		}
	}

	for _, currentItem := range current {
		found := false
		for _, desiredItem := range desired {
			if desiredItem == currentItem.String() {
				found = true
			}
		}
		if !found {
			deleteSlice = append(deleteSlice, currentItem.String())
		}
	}

	return addSlice, deleteSlice
}

func EnumValuesValidation(ty string, checkValuesSet interface{}, enumValues map[string]int32) diag.Diagnostics {
	values, ok := checkValuesSet.(*schema.Set)
	if !ok {
		return diag.Errorf("Attribute %s is no set for enum value check", ty)
	}

	for _, value := range values.List() {
		_, ok := enumValues[value.(string)]
		if !ok {
			return diag.Errorf("Attribute %s has unsupported enum value \"%s\"", ty, value)
		}
	}
	return nil
}

func EnumValueValidation(ty string, checkValue interface{}, enumValues map[string]int32) diag.Diagnostics {
	value, ok := checkValue.(string)
	if !ok {
		return diag.Errorf("Attribute %s is no string for enum value check", ty)
	}

	_, ok = enumValues[value]
	if !ok {
		return diag.Errorf("Attribute %s has unsupported enum value \"%s\"", ty, value)
	}
	return nil
}

func GetID(d *schema.ResourceData, idVar string) string {
	idStr := ""
	id, ok := d.GetOk(idVar)
	if ok {
		idStr = id.(string)
	} else {
		idStr = d.Id()
	}
	return idStr
}
