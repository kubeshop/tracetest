package installer

import "strings"

var kubernetes = installer{
	preChecks: []preChecker{
		kubectlChecker,
		helmChecker,
		localEnvironmentChecker,
	},
	configs:   []configurator{},
	installFn: kubernetesInstaller,
}

func kubernetesInstaller(config configuration, ui UI) {}

func helmChecker(ui UI) {
	if commandExists("helm") {
		ui.Println(ui.Green("✔ helm already installed"))
		return
	}

	ui.Warning("I didn't find helm in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install Helm", installHelm},
		{"Fix it manually", exitOption(
			"Check the helm install docs on https://helm.sh/docs/intro/install/",
		)},
	}, 0)

	option.fn(ui)

	if commandExists("helm") {
		ui.Println(ui.Green("✔ helm was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ helm could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func kubectlChecker(ui UI) {
	if commandExists("kubectl") {
		ui.Println(ui.Green("✔ kubectl already installed"))
		return
	}

	ui.Warning("I didn't find kubectl in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install Kubectl", installKubectl},
		{"Fix it manually", exitOption(
			"Check the kubectl install docs on https://kubernetes.io/docs/tasks/tools/#kubectl",
		)},
	}, 0)

	option.fn(ui)

	if commandExists("kubectl") {
		ui.Println(ui.Green("✔ kubectl was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ kubectl could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func localEnvironmentChecker(ui UI) {
	localK8sRunning := ui.Confirm("Do you have a local kubernentes running?", true)
	if !localK8sRunning {
		option := ui.Select("We can fix that:", []option{
			{"Install minikube", minikubeChecker},
			{"Fix manually", exitOption(
				"Check the minikube install docs on https://minikube.sigs.k8s.io/docs/start/",
			)},
		}, 0)

		option.fn(ui)
	}
}

func minikubeChecker(ui UI) {
	if commandExists("minikube") {
		ui.Println(ui.Green("✔ minikube already installed"))
		return
	}

	installMinikube(ui)

	if commandExists("minikube") {
		ui.Println(ui.Green("✔ minikube was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ minikube could not be installed. Check output for errors. " + createIssueMsg))
	}

	if !commandExists("docker") {
		ui.Warning("Docker is required to run your minikube instance")
		dockerChecker(ui)
	}
}

func installHelm(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing helm manually. See https://helm.sh/docs/intro/install/",
		installDocs:   "https://helm.sh/docs/intro/install/",
		args:          map[string]string{},
		apt: `
			sudo apt update -y
			sudo apt install -y curl gpg

			curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
			sudo apt-get install apt-transport-https --yes
			echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
			sudo apt-get update
			sudo apt-get install helm
		`,
		yum:      "sudo yum -y install helm",
		dnf:      "sudo dnf -y install helm",
		homebrew: "brew install helm",
		macIntelChipManual: `
			curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
			chmod 700 get_helm.sh
			./get_helm.sh
		`,
		macAppleChipManual: `
			curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
			chmod 700 get_helm.sh
			./get_helm.sh
		`,
		windows: "",
		other:   "",
	}).exec(ui)
}

func installKubectl(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing kubectl manually. See https://kubernetes.io/docs/tasks/tools/#kubectl",
		installDocs:   "https://kubernetes.io/docs/tasks/tools/#kubectl",
		args:          map[string]string{},
		apt: `
			sudo apt-get update
			sudo apt-get install -y ca-certificates curl apt-transport-https

			sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
			echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

			sudo apt-get update
			sudo apt-get install -y kubectl
		`,
		yum: `
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
sudo yum -y install kubectl
		`,
		dnf: `
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
sudo dnf -y install kubectl
		`,
		homebrew: "brew install helm",
		macIntelChipManual: `
			curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
			chmod +x ./kubectl
			sudo mv ./kubectl /usr/local/bin/kubectl
			sudo chown root: /usr/local/bin/kubectl
		`,
		macAppleChipManual: `
			curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/arm64/kubectl"
			chmod +x ./kubectl
			sudo mv ./kubectl /usr/local/bin/kubectl
			sudo chown root: /usr/local/bin/kubectl
		`,
		windows: "",
		other:   "",
	}).exec(ui)
}

func installMinikube(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing minikube manually. See https://minikube.sigs.k8s.io/docs/start/",
		installDocs:   "https://minikube.sigs.k8s.io/docs/start/",
		args: map[string]string{
			"Architecture":    detectArchitecture(),
			"RpmArchitecture": strings.Replace(detectArchitecture(), "amd64", "x86_64", 1),
		},
		apt: `
			sudo apt update
			sudo apt install curl -y
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_{{.Architecture}}.deb
			sudo dpkg -i minikube_latest_{{.Architecture}}.deb
		`,
		yum: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-latest.{{.RpmArchitecture}}.rpm
			sudo rpm -Uvh minikube-latest.{{.RpmArchitecture}}.rpm
		`,
		dnf: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-latest.{{.RpmArchitecture}}.rpm
			sudo rpm -Uvh minikube-latest.{{.RpmArchitecture}}.rpm
		`,
		homebrew: "brew install minikube",
		macIntelChipManual: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
			sudo install minikube-darwin-amd64 /usr/local/bin/minikube
		`,
		macAppleChipManual: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-arm64
			sudo install minikube-darwin-arm64 /usr/local/bin/minikube
		`,
		windows: "",
		other:   "",
	}).exec(ui)
}
