package diamond_client

import (
    "testing"
)

func TestSubstitute(t *testing.T) {
    config := "name=John\nfull=${this.name} Doe\nlong=${this.full} Richard\n"
    substitute, _ := Substitute(config, true, "", "", nil)
    if "name=John\nfull=John Doe\nlong=John Doe Richard\n" != substitute {
        t.Fail()
    }
}
