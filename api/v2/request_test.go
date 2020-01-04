package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBuilder(test *testing.T) {
	url := "https://api.telegram.org/bot989819653:AAFSXCy6cGafMK-VyBg9RVnehzrMRzrS9vM/sendMessage"
	var msg = map[string]interface{}{
		"text":    "Basic Message from Request Builder",
		"chat_id": 68386493,
	}
	req := NewBuilder(url).
		Post().
		AddHeader("Content-Type", "application/json").
		MarshalBody(msg).
		Build()
	ch := make(chan *ResponseChannel)
	go req.Run(ch)
	data := <-ch

	if data.Err != nil || data.Res == nil {
		test.Errorf("Error: %v, Res: %v", data.Err, data.Res)
	} else if data.Res.Status != 200 {
		test.Errorf("Res: %v", data.Res)
	} else {
		test.Logf("%v", data.Res)
	}
}

func TestKayWeird(test *testing.T) {
	i := 0
X:
	for n := 0; n < 5; n++ {
		println(n + 1)
		test.Log(n + 1)
	}
	if i < 5 {
		println("goto wtf")
		test.Log("goto wtf")
		i++
		goto X
	}
}

// func assert()
func TestBuilderImmutability(t *testing.T) {
	base := NewBuilder("one").Post()
	two := base.AppendUrl("/two")
	three := base.AppendUrl("/three").Build().(*request)
	multi := two.AppendUrl("/three").Build().(*request)
	twoB := two.Build().(*request)
	baseRe := base.Build().(*request)
	assert.Equal(t, "one", baseRe.url, "%v != %v", "one", baseRe.url)
	assert.Equal(t, "one/two", twoB.url, "%v != %v", "one", twoB.url)
	assert.Equal(t, "one/three", three.url, "%v != %v", "one", three.url)
	assert.Equal(t, "one/two/three", multi.url, "%v != %v", "one", three.url)
}
