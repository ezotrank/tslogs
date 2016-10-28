package tslogs

import (
	"testing"
	// "fmt"
)

func TestLoadConfigFileWithConfigV1(t *testing.T) {
	config, err := LoadConfigFile("./test/fixtures/configv1")
	if err != nil {
		t.Fatalf("can't load or parse config, err: %v", err)
	}
	if _, ok := config.Destinations["tsdb"]; !ok {
		t.Errorf("can't find destinations tsdb")
	}
	if _, ok := config.Destinations["datadogs"]; !ok {
		t.Errorf("can't find destinations tsdb")
	}
	if val, ok := config.Groups["app1"]; !ok {
		t.Errorf("can't load app1 in config")
	} else {
		if val.Mask != "/logs/app.log" {
			t.Errorf("mask not right, want /logs/app.log but got %v", val.Mask)
		}
		loadedDestinations := make(map[string]bool)
		for _,dst := range val.destinations {
			loadedDestinations[dst.Name()] = true
		}
		if len(loadedDestinations) != 2 {
			t.Fatalf("size of loaded destioantions should be 2 not %d", len(loadedDestinations))
		}
		for _, shouldLoadedDst := range []string{"tsdb", "datadogs"} {
			if _,ok := loadedDestinations[shouldLoadedDst]; !ok {
				t.Errorf("can't find destination %s in loaded", shouldLoadedDst)
			}
		}
		if val.Rules[0].Name != "app.handler_exec" {
			t.Errorf("rule name should be app.handler_exec not %q", val.Rules[0].Name)
		}
		if val.Rules[0].Regexp != `^\[.+\] (?P<http_code>\d+) (?P<http_method>GET|POST) .+ (?P<val>\d+\.\d+)ms$` {
			t.Errorf("rule regexp should not be %q", val.Rules[0].Regexp)
		}
		if val.Rules[0].Match != "] " {
			t.Errorf("rule regexp should not be %q", val.Rules[0].Match)
		}
	}

	if val, ok := config.Groups["app2"]; !ok {
		t.Errorf("can't load app2 in config")
	} else {
		if val.Mask != "/logs/app2.log" {
			t.Errorf("mask not right, want /logs/app2.log but got %v", val.Mask)
		}
		loadedDestinations := make(map[string]bool)
		for _,dst := range val.destinations {
			loadedDestinations[dst.Name()] = true
		}
		if len(loadedDestinations) != 2 {
			t.Fatalf("size of loaded destioantions should be 2 not %d", len(loadedDestinations))
		}
		for _, shouldLoadedDst := range []string{"tsdb", "datadogs"} {
			if _,ok := loadedDestinations[shouldLoadedDst]; !ok {
				t.Errorf("can't find destination %s in loaded", shouldLoadedDst)
			}
		}
		if val.Rules[0].Name != "app2.handler_exec" {
			t.Errorf("rule name should be app.handler_exec not %q", val.Rules[0].Name)
		}
		if val.Rules[0].Regexp != `^\[.+\] (?P<http_code>\d+) (?P<http_method>GET|POST) .+ (?P<val>\d+\.\d+)ms$` {
			t.Errorf("rule regexp should not be %q", val.Rules[0].Regexp)
		}
		if val.Rules[0].Match != "] " {
			t.Errorf("rule regexp should not be %q", val.Rules[0].Match)
		}
	}
}

// func TestLoadConfigFileWithConfigV2(t *testing.T) {
// 	_, err := LoadConfigFile("./test/fixtures/configv2")
// 	if err != nil {
// 		t.Fatalf("can't load or parse config, err: %v", err)
// 	}
// }