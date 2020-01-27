package pipeline

import (
	"fmt"
	"time"

	"github.com/project-flogo/catalystml-flogo/action/types"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/support/log"
)

type Instance struct {
	def *Definition
	id  string

	logger log.Logger
}

//  Get new instance from defination
func NewInstance(definition *Definition, id string, logger log.Logger) *Instance {

	return &Instance{def: definition, id: id, logger: logger}
}

func (inst *Instance) Id() string {
	return inst.id
}

// Run the instance of the CML.
func (inst *Instance) Run(input map[string]interface{}) (output map[string]interface{}, err error) {

	// Get the Scope of the CML pipeline.
	// Scope is the collection of the data in the CML
	scope, err := NewPipelineScope(input, inst.def.labels)

	if err != nil {
		return nil, err
	}

	// Log the time
	start := time.Now()

	//Check the type of the input of the pipeline.
	for key, _ := range inst.def.input {

		temp, ok := inst.def.input[key].(PipelineInput)
		if !ok {
			continue
		}

		err = types.ValidateType(temp.Type, input[key])

		if err != nil {
			return nil, err
		}

	}

	//Run the tasks.
	for _, task := range inst.def.tasks {

		scope, _ = task.Eval(scope, inst.logger)

	}

	// Set the output.

	if inst.def.output.Data != nil {
		mf := GetMapperFactory()
		mappings := make(map[string]interface{})

		// Type Switch
		switch t := inst.def.output.Data.(type) {
		case map[string]interface{}:
			for key, val := range t {
				mappings[key] = val
			}
		default:
			mappings["data"] = inst.def.output.Data
		}

		// Get the data from output expression
		outMapper, err := mf.NewMapper(mappings)
		output, err = outMapper.Apply(scope)

		if err != nil {
			return nil, err
		}
		var definedType data.Type

		// Check if the output is defined as dataframe or map.
		if inst.def.output.Type == "dataframe" || inst.def.output.Type == "map" {
			definedType, _ = data.ToTypeEnum("object")

			givenType, _ := data.GetType(output)

			if definedType != givenType {
				return nil, fmt.Errorf("Type mismatch in output. Defined type [%s] passed type [%s]", definedType, givenType)
			}

			inst.logger.Infof("The output took %v to calculate", time.Since(start))

			return output, nil
		}

		definedType, _ = data.ToTypeEnum(inst.def.output.Type)

		for key, _ := range output {

			givenType, _ := data.GetType(output[key])

			if definedType != givenType {
				return nil, fmt.Errorf("Type mismatch in output. Defined type [%s] passed type [%s]", definedType, givenType)
			}
		}

	}
	inst.logger.Infof("The output took %v to calculate", time.Since(start))

	return output, nil

}
