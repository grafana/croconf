package main

import (
	//"time"
	"go.k6.io/croconf"
	"go.k6.io/croconf/duration"
)

type ScriptConfig struct {
	cm *croconf.Manager
	*GlobalConfig

	UserAgent string
	VUs       int64

	Duration duration.Duration

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

	cm.AddField(
		croconf.NewStringField(
			&conf.UserAgent,
			croconf.DefaultStringValue("croconf example demo v0.0.1 (https://k6.io/)"),
			jsonSource.From("userAgent", nil),
			envVarsSource.From("K6_USER_AGENT", nil),
			cliSource.FromName("--user-agent"),
			// TODO: figure this out...
			// croconf.WithDescription("user agent for http requests")
		),
	)

	cm.AddField(croconf.NewInt64Field(
		&conf.VUs,
		croconf.DefaultInt64Value(1),
		jsonSource.From("vus", nil),
		envVarsSource.From("K6_VUS", nil),
		cliSource.FromNameAndShorthand("--vus", "-u"),
		// croconf.WithDescription("number of virtual users") // TODO
	))

	cm.AddField(croconf.NewCustomField(
		&conf.Duration,
		//jsonSource.From("duration", duration.FromJSON),
		envVarsSource.From("K6_DURATION", duration.FromEnv),
		//cliSource.FromNameAndShorthand("--duration", "-d"),
	))

	// TODO: add the other options and actually process and consolidate the
	// config values and handle any errors... Here we probably want to error out
	// if we see unknown CLI flags or JSON options

	if err := cm.Consolidate(); err != nil {
		return nil, err
	}

	return conf, nil
}
