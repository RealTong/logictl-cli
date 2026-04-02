package config

type Config struct {
	Daemon   DaemonConfig `toml:"daemon"`
	Devices  []Device     `toml:"devices"`
	Actions  []Action     `toml:"actions"`
	Profiles []Profile    `toml:"profiles"`
}

type DaemonConfig struct {
	Enabled bool `toml:"enabled"`
}

type Device struct {
	Name string `toml:"name"`
}

type Action struct {
	ID      string `toml:"id"`
	Type    string `toml:"type"`
	Command string `toml:"command,omitempty"`
}

type Profile struct {
	Name     string    `toml:"name"`
	App      string    `toml:"app"`
	Bindings []Binding `toml:"bindings"`
}

type Binding struct {
	Trigger string `toml:"trigger"`
	Action  string `toml:"action"`
}
