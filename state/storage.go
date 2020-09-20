package state

//var templates = make([]Template, 0, 10)
var templates = make(map[int]Template)
var staging = make(map[int]State)

func NewTemplate(targetChannel string, srcPtr MessagePtr, text string) Template {
	maxID := 0
	for _, t := range templates {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1
	return Template{ID: newID, TargetChannel: targetChannel, SourceMessagePtr: srcPtr, Text: text}
}

func GetTemplateBySource(srcPtr MessagePtr) (tpl Template, ok bool) {
	for k, t := range templates {
		if t.SourceMessagePtr == srcPtr {
			return templates[k], true
		}
	}
	return Template{}, false
}

func GetTemplateByID(id int) (tpl Template, ok bool) {
	for i, t := range templates {
		if t.ID == id {
			return templates[i], true
		}
	}
	return Template{}, false
}

func GetTemplatesCount() int {
	return len(templates)
}

func GetTemplates() []Template {
	tpls := make([]Template, len(templates))
	i := 0
	for _, v := range templates {
		tpls[i] = v
		i++
	}
	return tpls
}

func GetTemplatesWithState(state State) []Template {
	tpls := make([]Template, 0, len(staging))
	for id, s := range staging {
		if s != state {
			continue
		}
		tpl, _ := templates[id]
		tpl.ID = id // For deleted templates
		tpls = append(tpls, tpl)
	}
	return tpls
}

func setStateForID(id int, state State) {
	if state == Null {
		delete(staging, id)
		return
	}
	if state == Clean {
		staging[id] = Clean
		return
	}
	var oldState, newState State
	if checkState, ok := staging[id]; ok {
		oldState = checkState
	} else {
		oldState = Null
	}
	if oldState == state {
		return
	}
	switch oldState {
	case Null:
		newState = state
	case Added:
		if state == Deleted {
			newState = Null
		} else {
			newState = oldState
		}
	case Updated, Deleted, Clean: // Pre-existente state
		if state == Added {
			newState = Updated
		} else {
			newState = state
		}
	}
	staging[id] = newState
}

func SetTemplate(tpl Template) {
	var state State
	if _, ok := templates[tpl.ID]; ok {
		state = Updated
	} else {
		state = Added
	}
	templates[tpl.ID] = tpl
	setStateForID(tpl.ID, state)
}

func DeleteTemplateByID(id int) int {
	was := len(templates)
	delete(templates, id)
	became := len(templates)
	if was-became > 0 {
		setStateForID(id, Deleted)
		PersistTemplates()
	}
	return was - became
}
