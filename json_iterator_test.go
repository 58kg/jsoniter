package jsoniter

import (
	"fmt"
	"testing"

	"github.com/58kg/to_string"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTraverse(test *testing.T) {
	Convey("TestTraverse", test, func() {
		Convey("array", func() {
			const inputStr = `  [
								  {
									"k1": 1,
									"k2": {
									  "k3": "v3"
									}
								  },
								  2,
								  3,
								  4,
								  {
									"k5": [
									  {
										"k4": "v4"
									  },
									  {
										"k4": "v4"
									  }
									]
								  }
								]`
			var arr []interface{}
			err := json.Unmarshal([]byte(inputStr), &arr)
			So(err, ShouldEqual, nil)
			err = Traverse(arr, []Handler{
				{
					Fields: []string{"k1"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k1")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, 1)
						obj.(map[string]interface{})[fields[len(fields)-1]] = 2
						return nil
					},
				},
				{
					Fields: []string{"k2", "k3"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k3")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, "v3")
						obj.(map[string]interface{})[fields[len(fields)-1]] = "v2"
						return nil
					},
				},
				{
					Fields: []string{"k5", "k4"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k4")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, "v4")
						obj.(map[string]interface{})[fields[len(fields)-1]] = "v5"
						return nil
					},
				},
				{
					Fields: []string{"k5"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k5")
						So(len(obj.(map[string]interface{})[fields[len(fields)-1]].([]interface{})), ShouldEqual, 2)
						return nil
					},
				},
			})
			So(err, ShouldEqual, nil)
			str, err := json.Marshal(arr)
			So(err, ShouldEqual, nil)
			const outputStr = `[{"k1":2,"k2":{"k3":"v2"}},2,3,4,{"k5":[{"k4":"v5"},{"k4":"v5"}]}]`
			So(string(str), ShouldEqual, outputStr)
		})

		Convey("map", func() {
			const inputStr = `{
					  "k1": 1,
					  "k2": 0.5,
					  "k3": [
						1,
						2,
						3,
						{
						  "k4": "v1"
						},
						[
						  4,
						  5,
						  6,
						  {
							"k4": "v1"
						  }
						]
					  ],
					  "k5": {
						"k6": {
						  "k7": [
							{
							  "k8": "v8"
							},
							8,
							9
						  ]
						}
					  }
					}`
			var m map[string]interface{}
			err := json.Unmarshal([]byte(inputStr), &m)
			So(err, ShouldEqual, nil)
			err = Traverse(m, []Handler{
				{
					Fields: []string{"k1"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k1")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, 1)
						obj.(map[string]interface{})[fields[len(fields)-1]] = 2
						return nil
					},
				},
				{
					Fields: []string{"k3", "k4"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k4")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, "v1")
						obj.(map[string]interface{})[fields[len(fields)-1]] = "v2"
						return nil
					},
				},
				{
					Fields: []string{"k5", "k6", "k7"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k7")
						So(obj.(map[string]interface{})[fields[len(fields)-1]].([]interface{})[2], ShouldEqual, 9)
						obj.(map[string]interface{})[fields[len(fields)-1]].([]interface{})[2] = 10
						return nil
					},
				},
				{
					Fields: []string{"k5", "k6", "k7", "k8"},
					Handler: func(obj interface{}, fields []string) error {
						fmt.Printf("fields:%v, obj:%v\n", to_string.String(fields), to_string.StringByConf(obj, to_string.Config{
							InformationLevel: to_string.AllTypesInfo,
						}))
						So(fields[len(fields)-1], ShouldEqual, "k8")
						So(obj.(map[string]interface{})[fields[len(fields)-1]], ShouldEqual, "v8")
						obj.(map[string]interface{})[fields[len(fields)-1]] = "v9"
						return nil
					},
				},
			})
			So(err, ShouldEqual, nil)
			str, err := json.Marshal(m)
			So(err, ShouldEqual, nil)
			const outputStr = `{"k1":2,"k2":0.5,"k3":[1,2,3,{"k4":"v2"},[4,5,6,{"k4":"v2"}]],"k5":{"k6":{"k7":[{"k8":"v9"},8,10]}}}`
			So(string(str), ShouldEqual, outputStr)
		})
	})
}
