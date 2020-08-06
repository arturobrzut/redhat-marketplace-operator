package razeedeployment_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestRazeedeployment(t *testing.T) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Razeedeployment Suite")
}
