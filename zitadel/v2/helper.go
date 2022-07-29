package v2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type stringified struct {
	str string
}

func (s *stringified) String() string {
	return s.str
}

type stringify interface {
	String() string
}

func setToStringSlice(set *schema.Set) []string {
	slice := make([]string, 0)
	for _, secondFactor := range set.List() {
		slice = append(slice, secondFactor.(string))
	}
	return slice
}

func getAddAndDelete(current []stringify, desired []string) ([]string, []string) {
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
