package plugins

import (
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
    p.Register()
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
    return "Ok."
}

// We check if the function is in our precompiled list of plugins then use reflect to 
// be able to call the corresponding method via a string call. Since we have to export the 
// methods, we change the first char of string s to uppercase.
func (p *Plugins) Call(s []string) []string {
    u := []rune(s[0])
    u[0] = unicode.ToUpper(u[0])
    s[0] = string(u)
    for _, e := range p.List {
        if s[0] == e {
            in := make([]reflect.Value, 1)
            in[0] = reflect.ValueOf(s)
            methodValue := reflect.ValueOf(p).MethodByName(s[0]).Call(in)
            var ret []string
            for i := 0; i < methodValue[0].Len(); i++ {
                ret = append(ret, methodValue[0].Index(i).String())
            }
            return ret
        }
    }
    return []string{}
}
