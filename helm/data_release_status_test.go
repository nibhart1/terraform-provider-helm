package helm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccStatusRelease_basic(t *testing.T) {
	name := fmt.Sprintf("test-basic-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("%s-%s", testNamespace, acctest.RandString(10))
	// Delete namespace automatically created by helm after checks
	defer deleteNamespace(t, namespace)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHelmReleaseDestroy(namespace),
		Steps: []resource.TestStep{{
			Config: testAccHelmStatusConfigBasic(name, namespace, "0.6.2"),

			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("helm_release.test", "metadata.0.name", name),
				resource.TestCheckResourceAttr("helm_release.test", "metadata.0.namespace", namespace),
				resource.TestCheckResourceAttr("helm_release.test", "metadata.0.revision", "1"),
				resource.TestCheckResourceAttr("helm_release.test", "status", "DEPLOYED"),
				resource.TestCheckResourceAttr("helm_release.test", "metadata.0.chart", "mariadb"),
				resource.TestCheckResourceAttr("helm_release.test", "metadata.0.version", "0.6.2"),
				resource.TestCheckResourceAttr("data.helm_release_status.test", "pod.#", "1"),
				resource.TestCheckResourceAttr("data.helm_release_status.test", "service.#", "1"),
				resource.TestCheckResourceAttr("data.helm_release_status.test", "pvc.#", "1"),
				resource.TestCheckResourceAttr("data.helm_release_status.test", "ingress.#", "1"),
			),
		}},
	})
}

func testAccHelmStatusConfigBasic(name, namespace, version string) string {
	return fmt.Sprintf(`
		resource "helm_release" "test" {
 			name      = "%s"
			namespace = "%s"
  			chart     = "stable/mariadb"
			version   = "%s"

			set {
				name = "foo"
				value = "qux"
			}

			set {
				name = "qux.bar"
				value = 1
			}
			wait = true
		}
		data "helm_release_status" "test" {
 			name      = "${helm_release.test.metadata.0.name}"
			revision   = 1


		}
	`, name, namespace, version)
}
