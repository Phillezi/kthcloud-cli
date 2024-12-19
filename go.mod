module github.com/Phillezi/kthcloud-cli

go 1.23.2

require (
	github.com/briandowns/spinner v1.23.1
	github.com/browserutils/kooky v0.2.2
	github.com/go-resty/resty/v2 v2.15.3
	github.com/google/uuid v1.6.0
	github.com/jedib0t/go-pretty/v6 v6.6.1
	github.com/manifoldco/promptui v0.9.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	go-deploy v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Velocidex/json v0.0.0-20220224052537-92f3c0326e5a // indirect
	github.com/Velocidex/ordereddict v0.0.0-20230909174157-2aa49cc5d11d // indirect
	github.com/Velocidex/yaml/v2 v2.2.8 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/analysis v0.22.0 // indirect
	github.com/go-openapi/errors v0.21.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/loads v0.21.4 // indirect
	github.com/go-openapi/spec v0.20.13 // indirect
	github.com/go-openapi/strfmt v0.21.10 // indirect
	github.com/go-openapi/swag v0.22.6 // indirect
	github.com/go-openapi/validate v0.22.4 // indirect
	github.com/go-sqlite/sqlite3 v0.0.0-20180313105335-53dd8e640ee7 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gonuts/binary v0.2.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/keybase/go-keychain v0.0.0-20231219164618-57a3676c3af6 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/openshift/api v0.0.0-20230503133300-8bbcb7ca7183 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/zalando/go-keyring v0.2.6 // indirect
	go.mongodb.org/mongo-driver v1.13.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.29.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.3 // indirect
	k8s.io/apimachinery v0.29.1 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	kubevirt.io/api v1.2.0 // indirect
	kubevirt.io/containerized-data-importer-api v1.57.0-alpha1 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	www.velocidex.com/golang/go-ese v0.2.0 // indirect
)

replace go-deploy => github.com/kthcloud/go-deploy v0.0.0-20240613154044-4662d95a90ad

replace github.com/browserutils/kooky v0.2.2 => github.com/DP19/kooky v0.2.4
