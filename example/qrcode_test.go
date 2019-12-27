package example

import "testing"

func TestQRCode(t *testing.T) {
	err := NewQRCode("http://www.baidu.com", "./baidu.png")
	if err != nil {
		t.Error(err)
	}
}
