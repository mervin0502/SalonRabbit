// Attribute package
//
package attribute

import "mervin.me/util"

// Attribute:
// basic attribute:1,label;2,weight;3,
// other attribute:
//
//
//register the element attribute before create the element object.
type Attribute struct {
	value map[string]interface{} "the attribte value"
	size  int                    "the number of attribute in the element"
}

func New() *Attribute {
	return &Attribute{make(map[string]interface{}), 0}
}

// the attribute is empty
func (a *Attribute) IsEmpty() bool {
	if a.size != 0 {
		return true
	} else {
		return false
	}
}

// the attribute size
func (a *Attribute) Size() int {
	return a.size
}

//clear the attribute
//a.value =  make(map[string]interface {},0)?
func (a *Attribute) Clear() {
	a.value = make(map[string]interface{}, 0)
	a.size = 0
}

//copy a new attribute
func (a *Attribute) Copy() *Attribute {
	newValue := make(map[string]interface{}, a.size)
	for k, v := range a.value {
		newValue[k] = v
	}
	ac := &Attribute{value: newValue, size: a.size}
	return ac
}

func (a *Attribute) String() string {
	var str string = "attribute:{"
	for k, v := range a.value {
		str += k + ":" + util.String(v) + " "
	}
	str += "}"
	return str
}

//PutAttribute attribute
func (a *Attribute) Put(key string, value interface{}) {

	if _, ok := a.value[key]; !ok {
		a.value[key] = value
		a.size++
	} else {
		a.value[key] = value
	}
}

//get attribute
func (a *Attribute) Get(k string) interface{} {
	if v, ok := a.value[k]; ok {
		return v
	} else {
		return nil
	}

}

//get attribute key
func (a *Attribute) GetKeys() []string {
	keys := make([]string, a.size)
	i := 0
	for k, _ := range a.value {
		keys[i] = k
		i++
	}
	return keys
}

// get the all attribute
func (a *Attribute) GetAll() map[string]interface{} {
	m := make(map[string]interface{}, a.size)
	for k, v := range a.value {
		m[k] = v
	}
	return m
}

//remove a attribute of element
func (a *Attribute) Delete(k string) {
	delete(a.value, k)
}

//has attribute k
func (a *Attribute) Has(k string) bool {
	_, ok := a.value[k]
	return ok
}
