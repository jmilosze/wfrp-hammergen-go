package warhammer

import (
	"fmt"
)

func GetValidationAliases() map[string]string {
	alias := make(map[string]string)

	alias["name_valid"] = "min=0,max=200,excludesall=<>"
	alias["desc_valid"] = "min=0,max=100000,excludesall=<>"
	alias["shared_valid"] = "boolean"
	alias["medium_string_valid"] = "min=0,max=200,excludesall=<>"
	alias["source_valid"] = fmt.Sprintf("dive,keys,oneof=%s,endkeys,min=0,max=15,excludesall=<>", getAllowedSourceType())

	alias["mutation_type_valid"] = fmt.Sprintf("oneof=%s", getAllowedMutationType())
	alias["property_type_valid"] = fmt.Sprintf("oneof=%s", getAllowedPropertyType())
	alias["property_applicable_to_valid"] = fmt.Sprintf("dive,oneof=%s", getAllowedItemType())

	return alias
}
