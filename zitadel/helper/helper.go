package helper

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func GetOkSetToStringSlice(d *schema.ResourceData, value string) []string {
	var slice []string
	if set, ok := d.GetOk(value); ok {
		slice = SetToStringSlice(set.(*schema.Set))
	} else {
		slice = make([]string, 0)
	}
	return slice
}

func SetToStringSlice(set *schema.Set) []string {
	slice := make([]string, 0)
	for _, secondFactor := range set.List() {
		slice = append(slice, secondFactor.(string))
	}
	return slice
}

func GetAddAndDelete(current []string, desired []string) ([]string, []string) {
	addSlice := make([]string, 0)
	deleteSlice := make([]string, 0)

	for _, desiredItem := range desired {
		found := false
		for _, currentItem := range current {
			if desiredItem == currentItem {
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
			if desiredItem == currentItem {
				found = true
			}
		}
		if !found {
			deleteSlice = append(deleteSlice, currentItem)
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

func GetStringFromAttr(ctx context.Context, attrs map[string]attr.Value, key string) string {
	value, err := attrs[key].ToTerraformValue(ctx)
	if err != nil {
		return ""
	}
	var str string
	if err := value.As(&str); err != nil {
		return ""
	}
	return str
}

func DescriptionEnumValuesList(enum map[int32]string) string {
	str := ", supported values: "
	values := make([]string, len(enum))
	highest := 0
	for k := range enum {
		if int(k) > highest {
			highest = int(k)
		}
	}

	j := 0
	for i := 0; i < highest+1; i++ {
		if value, ok := enum[int32(i)]; ok {
			values[j] = value
			j++
		}
	}
	str += strings.Join(values, ", ")
	return str
}

func EnumValueMap(enum map[int32]string) map[string]int32 {
	values := make(map[string]int32)
	for k, v := range enum {
		values[v] = k
	}
	return values
}
