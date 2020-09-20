package services

var variables = make(map[string]Variable)

type Variable struct {
	Name        string
	Value       string
	Description string
}

func RegisterVariable(name string, description string) {
	variables[name] = Variable{Name: name, Description: description}
}

func GetVariablesInfo() map[string]string {
	info := make(map[string]string)
	for k, v := range variables {
		info[k] = v.Description
	}
	return info
}

func GetValue(k string) string {
	if v, ok := variables[k]; ok {
		return v.Value
	} else {
		return ""
	}
}

func SetValue(key string, val string) bool {
	if v, ok := variables[key]; ok {
		v.Value = val
		variables[key] = v
		return true
	} else {
		return false
	}
}
