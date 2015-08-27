package plugins

import (
    "fmt"
    "strings"
    "unicode"
    "reflect"
)

type Plugins struct {
    List  []string
    InternalFuncs []string
}

// Plugins constructor, constructs a list of internal funcs
// not intended for the bot to list.
func NewPlugins() *Plugins {
    p := new(Plugins)
    p.InternalFuncs = []string{"Register", "Call"}
    return p
}

// Register all implemented methods of plugins package.
// The function excludes some plugin internal functions.
func (p *Plugins) Register() string {
    pType := reflect.TypeOf(p)
    for i := 0; i < pType.NumMethod(); i++ {
        method := pType.Method(i)
        p.List = append(p.List, method.Name)
    }
    for _, f := range p.InternalFuncs {
        for i, e := range p.List {
            if f == e {
                p.List = append(p.List[:i], p.List[i+1:]...)
            }
        }
    }
    fmt.Println(p.List)
    return "Ok."
}

// We check if th function is in our precompiled list of plugins then use reflect to 
// be able to call the corresponding method via a string call. Since we have to export the 
// methods, we change the first char of string s to uppercase.
func (p *Plugins) Call(s string) string {
    u := []rune(s)
    u[0] = unicode.ToUpper(u[0])
    s = string(u)
    for _, e := range p.List {
        if s == e {
            methodValue := reflect.ValueOf(p).MethodByName(s).Call([]reflect.Value{})
            var ret []string
            for _, v := range methodValue {
                ret = append(ret, v.String())
            }
            return strings.Join(ret, " ")
        }
    }
    return ""
}
