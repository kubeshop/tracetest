package installer

func Start() {
	DockerCompose.PreCheck(DefaultUI)
}

type Installer struct {
	preChecks []preChecker
}

func (i Installer) PreCheck(ui UI) {
	for _, pc := range i.preChecks {
		pc(ui)
	}
}

func (i Installer) Install() {}

type preChecker func(ui UI)
