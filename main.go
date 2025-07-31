package main

import (
	"fmt"
	"net/http"
	"reflect"
)

type AuthHeader struct {
	Header string `header:"Authorization"`
}

type Adapter struct {
	fn any
}

func Decorate(fn any) *Adapter {
	return &Adapter{fn: fn}
}

func (a *Adapter) Respond(resp http.ResponseWriter, req *http.Request) {
		fn := reflect.ValueOf(a.fn)
		if fn.Kind() == reflect.Func {
			fmt.Println("it's a function!")
		}

		fnTyp := fn.Type()
		args := make([]reflect.Value, fnTyp.NumIn())
		for inI := 0; inI < fnTyp.NumIn(); inI++ {
			argType := fnTyp.In(inI)
			if argType.Kind() == reflect.Pointer {
				args[inI] = reflect.New(argType)
				hydrate(args[inI], req)
			} else {
				arg := reflect.New(argType)
				hydrate(arg, req)
				args[inI] = reflect.Indirect(arg)
			}
		}

		fn.Call(args)

		fmt.Println(args)

		fmt.Println(fnTyp)
}


// hydrate fills in the reflect.Value for a given http argument by reading the request.
// for example, if the arg was a zerod out JSON object, we would parse the JSON object from
// the request.Body()
func hydrate(emptyArg reflect.Value, req *http.Request) {
	fmt.Println("emptyArg", emptyArg)
	argType := emptyArg.Elem()
	for i := 0; i < argType.NumField(); i++ {
		fieldType := argType.Field(i)
		reflect.Indirect(emptyArg).Field(i).SetString("foo")
		// todo get struct tag and set based on
	}
}

type MyArgs struct {
	X int `json:"x"`
}


// return the type of struct tags attached to a given reflect.Type (e.g. json, queryparam, form, etc.)
func getStructTypes(t reflect.Type) map[string]struct{} {

	tags := make(map[string]struct{})
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		fieldStrucTag := t.Field(i).Tag

		if _, ok := fieldStrucTag.Lookup("json"); ok {
				tags["json"] = struct{}{}
		}
		if _, ok := fieldStrucTag.Lookup("header"); ok {
				tags["header"] = struct{}{}
		}
	}

	return tags
}

func myendpoint(auth AuthHeader) map[string]string {
	fmt.Println("in handler: ", auth)
	return map[string]string{"foo": "bar"}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", Decorate(myendpoint).Respond)

	srv := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fmt.Println("listening...")
	srv.ListenAndServe()

	fmt.Println("hello")
}
