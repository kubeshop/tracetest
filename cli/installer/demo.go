package installer

func ConfigureDemoApp(conf configuration, ui UI) configuration {
	conf["demo.enable"] = ui.Confirm("Do you want to enable the demo app?", true)
	return conf
}
