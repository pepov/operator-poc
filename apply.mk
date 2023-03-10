apply:
	#go install k8s.io/code-generator/cmd/openapi-gen@latest
	GOPATH="" openapi-gen -p api/v1beta1/openapi -h hack/boilerplate.go.txt -i github.com/pepov/operator-poc/api/v1beta1 --input-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version"
	go run hack/crd-schema/main.go > api/v1beta1/openapi/openapi_generated.json
	#go install k8s.io/code-generator/cmd/applyconfiguration-gen@latest
	applyconfiguration-gen -i github.com/pepov/operator-poc/api/v1beta1 -p github.com/pepov/operator-poc/api/v1beta1/applyconfigurations -h hack/boilerplate.go.txt --openapi-schema api/v1beta1/openapi/openapi_generated.json --trim-path-prefix github.com/pepov/operator-poc