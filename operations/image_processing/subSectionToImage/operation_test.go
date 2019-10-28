package subsectiontoimage

import (
	"fmt"
	"io/ioutil"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/project-flogo/catalystml-flogo/action/support/test"
	"github.com/stretchr/testify/assert"
)

func TestFalse(t *testing.T) {

	inputs := make(map[string]interface{})

	p := Params{Size: []int{1500, 1750}, LowerLeftCorner: []int{3000, 2000}}

	file := "../test_image.jpg"
	img, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error Oopening file: %v\n", err)
		return
	}

	inputs["img"] = img

	optInitConext := test.NewOperationInitContext(p, nil)
	opt, err := New(optInitConext)
	assert.Nil(t, err)

	output, err := opt.Eval(inputs)
	assert.NotNil(t, output)
	assert.Nil(t, err)

}
