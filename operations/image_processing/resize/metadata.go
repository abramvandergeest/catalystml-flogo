package resize

import (
	"image"
	// "github.com/project-flogo/core/data/coerce"
)

type Params struct {
	Xsize int    `md:"xsize",required=false`
	Ysize int    `md:"ysize",required=false`
	Algo  string `md:"algo",allowed=["Lanczos","NearestNeighbor","Linear", "CatmullRom"],required=false`
}

type Config struct {
	Operation string                 `json:"operation"`
	Params    map[string]interface{} `json:"params,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
}

type Input struct {
	Img image.Image `md:"img"`
}

func (i *Input) FromMap(values map[string]interface{}) error {

	i.Img = values["img"].(image.Image)

	return nil
}
