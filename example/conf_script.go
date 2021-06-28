package main

import "go.k6.io/croconf"

type ScriptConfig struct {
	*GlobalConfig
	cm *croconf.Manager

	UserAgent string
	VUs       int64

	// TODO: have a sub-config
}

func NewScriptConfig(
	globalConf *GlobalConfig,
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
	jsonSource *croconf.SourceJSON,
) (*ScriptConfig, error) {
	cm := croconf.NewManager()
	conf := &ScriptConfig{GlobalConfig: globalConf, cm: cm} // TODO: somehow save the sources in the struct as well?

	cm.StringField(
		&conf.UserAgent,
		croconf.DefaultStringValue("croconf example demo v0.0.1 (https://k6.io/)"),
		jsonSource.From("userAgent"),
		envVarsSource.From("K6_USER_AGENT"),
		cliSource.FromName("--user-agent"),
		// TODO: figure this out...
		// croconf.WithDescription("user agent for http requests")
	)

	cm.Int64Field(
		&conf.VUs,
		croconf.DefaultInt64Value(1),
		jsonSource.From("vus"),
		envVarsSource.From("K6_VUS"),
		cliSource.FromNameAndShorthand("--vus", "-u"),
		// croconf.WithDescription("number of virtual users") // TODO
	)

	// TODO: add the other options and actually process and consolidate the
	// config values and handle any errors... Here we probably want to error out
	// if we see unknown CLI flags or JSON options

	return conf, nil
}