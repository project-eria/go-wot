package producer

import (
	"sync"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/thing"
)

// Regression test for audit finding #4: concurrent GetPropertyChangeChannel /
// GetEventChannel while Emit* iterates the channel slices. Run with -race.
func Test_ChannelSlicesConcurrency(t *testing.T) {
	td, err := thing.New("urn:test:race", "1.0.0", "Race", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	boolData, _ := dataSchema.NewBoolean()
	prop := interaction.NewProperty("on", "On", "", boolData, interaction.PropertyObservable(true))
	td.AddProperty(prop)
	event := interaction.NewEvent("ding", "Ding", "")
	td.AddEvent(event)

	exposed := NewExposedThing(td, "race", &sync.WaitGroup{})
	exposed.SetEventHandler("ding", func() (interface{}, error) { return true, nil })

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			ch := exposed.GetPropertyChangeChannel()
			go func() { // drain so Emit never logs "channel blocked"
				for range ch {
				}
			}()
		}()
		go func() {
			defer wg.Done()
			exposed.GetEventChannel()
		}()
		go func() {
			defer wg.Done()
			_ = exposed.EmitPropertyChange("on", true, nil)
			_ = exposed.EmitEvent("ding", nil)
		}()
	}
	wg.Wait()
}
