package config

type Config struct {
	Devices  []Device     `toml:"devices"`
	Actions  []Action     `toml:"actions"`
	Profiles []Profile    `toml:"profiles"`
}

type Device struct {
	ID             string            `toml:"id"`
	MatchVendorID  int               `toml:"match_vendor_id"`
	MatchProductID int               `toml:"match_product_id"`
	Scroll         ScrollConfig      `toml:"scroll"`
	Capabilities   map[string]string `toml:"capabilities"`
}

type ScrollConfig struct {
	Direction    string `toml:"direction"`
	SmoothScroll bool   `toml:"smooth_scroll"`
}

type Action struct {
	ID     string   `toml:"id"`
	Type   string   `toml:"type"`
	Keys   []string `toml:"keys,omitempty"`
	System string   `toml:"system,omitempty"`
	Script string   `toml:"script,omitempty"`
}

type Profile struct {
	ID          string    `toml:"id"`
	AppBundleID string    `toml:"app_bundle_id"`
	Bindings    []Binding `toml:"bindings"`
}

type Binding struct {
	Device   string `toml:"device,omitempty"`
	Trigger  string `toml:"trigger"`
	Action   string `toml:"action"`
	Priority *int   `toml:"priority,omitempty"`
}
