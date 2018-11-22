package producer

import (
	"github.com/bgadrian/fastfaker/faker"
	//nolint
)

// GetRandomUserProperties creates a map with predefined user properties
func GetRandomUserProperties() map[string]string {
	//it uses the fastFaker template system
	//for the supported {variable} see https://github.com/bgadrian/fastfaker/blob/master/TEMPLATE_VARIABLES.md
	result := map[string]string{
		//Update the README.md if you change the keys
		KeyUserID:   "uuid",
		"sessionId": "uuid",
		"userAgent": "useragent",
		"country":   "country",
	}

	//get random data
	fastFaker := faker.NewFastFaker()
	for k, v := range result {
		result[k], _ = fastFaker.TemplateCustom(v, "", "")
	}

	return result
}

// NewRandomEventGenerator creates a Fake event generator based on a predefined set of patterns
// templateEvents must contains properties, and in their values FastFaker templates.
// ["price" => "###.##", "productCategory" => "{word}]
// each event generated will be based on these patterns, with fake data replacing the templates
func NewRandomEventGenerator(templateEvents []Event) EventGenerator {
	gen := &randomEventGenerator{}
	gen.templateEvents = templateEvents
	return gen
}

type randomEventGenerator struct {
	templateEvents []Event
}

func (g *randomEventGenerator) NewEvents(count int) []Event {
	//for now we keep our logic simple, with no correlation between the events
	//we just select them at random

	result := make([]Event, 0, count)
	fastFaker := faker.NewFastFaker()

	for len(result) < count {
		patternIndex := fastFaker.Intn(len(g.templateEvents))
		patternEvent := g.templateEvents[patternIndex]
		patternProps := patternEvent.Properties()

		props := make(map[string]string, len(patternProps))
		for key, template := range patternProps {
			props[key] = fastFaker.Template(template)
		}

		result = append(result, NewSimpleEvent(patternEvent.Name(), props))
	}
	return result
}
