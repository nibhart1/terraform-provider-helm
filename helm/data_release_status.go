package helm

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

var Totalpodscount int
var TotalServiceCount int
var TotalPvcCount int
var TotalIngressCount int

func dataSourceReleaseStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataStatusRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Release name.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "revision of the release.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "namespace of release",
			},
			"chart_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status of release",
			},

			"total_pods": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total number of pods",
			},
			"total_services": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total number of service",
			},
			"total_pvcs": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total number of pvc",
			},
			"total_ingress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total number of pvc",
			},

			"pod": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Pods deployed in that release",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the pod",
						},
						"age": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Age of the pod",
						},
					},
				},
			},
			"service": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service deployed in that release",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the service",
						},
						"age": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "age of the service being created",
						},
					},
				},
			},
			"pvc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Pvc resource in that release",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the pvc",
						},
						"age": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "age of the pvc being created",
						},
					},
				},
			},
			"ingress": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Ingress resource in that release",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the ingress",
						},
						"age": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "age of the ingress being created",
						},
					},
				},
			},
		},
	}
}

func dataStatusRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	client, err := m.GetHelmClient()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	revision := int32(d.Get("revision").(int))

	res, err := client.ReleaseStatus(name, helm.StatusReleaseVersion(revision)) //getting helm status of  deployed chart
	if err != nil {
		return err
	}
	Status := res.Info.Status.Code.String()
	Namespace := res.GetNamespace()
	d.Set("chart_status", Status)
	d.Set("namespace", Namespace)
	if len(res.Info.Status.Resources) > 0 {

		d.Set("pod", podStatusdetails(res, d))
		d.Set("total_pods", Totalpodscount)
		d.Set("service", serviceStatusdetails(res, d))
		d.Set("total_services", TotalServiceCount)
		d.Set("pvc", pvcStatusdetails(res, d))
		d.Set("total_pvcs", TotalPvcCount)
		d.Set("ingress", ingressStatusdetails(res, d))
		d.Set("total_ingress", TotalIngressCount)

	}
	d.SetId(name)
	return nil

}

func podStatusdetails(res *rls.GetReleaseStatusResponse, d *schema.ResourceData) []map[string]interface{} {
	var pods []string // getting all the pods available in the chart
	pods = make([]string, 0)
	var podcount = 0 //for counting the no of pods fetched

	pattern := `==> v\d/Pod\W*\w*\W*\nNAME\s+AGE\n`

	resources := strings.Split(res.Info.Status.Resources, "\n\n") //splitting the resource string on \n\n
	for _, res := range resources {

		matched, _ := regexp.MatchString(pattern, res) //fetching string containing info of pods
		if matched {
			pods = append(pods, res)
			podcount++
		}

	}
	Totalpodscount = podcount
	var totalPods []map[string]interface{}
	//creating a list of pods

	for pd := 0; pd < podcount; pd++ {

		/*format of each received pods as a single string
		  ==> v1/Pod
		  NAME         AGE
		  myapp-pod    3d22h

		*/

		resource := strings.Split(pods[pd], "AGE\n") //splitting after AGE
		for in, res2 := range resource {

			if in == 1 {
				morepods := strings.Split(res2, "\n")
				for _, res3 := range morepods {
					tempPod := make(map[string]interface{})

					//replacing all the extras whitespace between the strings
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(res3, " ")

					//spltting the pods string on whitespace
					splitpods := strings.Split(s, " ")

					tempPod["name"] = splitpods[0]
					tempPod["age"] = splitpods[1]

					totalPods = append(totalPods, tempPod)

				}

			}
		}

	}
	return totalPods
}

func serviceStatusdetails(res *rls.GetReleaseStatusResponse, d *schema.ResourceData) []map[string]interface{} {

	var services []string // get the service string available in the chart
	services = make([]string, 0)
	var servicecount = 0 //for counting the no of services fetched

	pattern := `==> v\d/Service\W*\w*\W*\nNAME\s+AGE\n`

	resources := strings.Split(res.Info.Status.Resources, "\n\n") //splitting the resource string on \n\n
	for _, res := range resources {

		matched, _ := regexp.MatchString(pattern, res) //fetching string containing info of services
		if matched {
			//services = res
			services = append(services, res)
			servicecount++
		}

	}

	TotalServiceCount = servicecount
	var totalServices []map[string]interface{}
	//creating a list of services

	for sv := 0; sv < servicecount; sv++ {
		/*format of each received service as a single string
					  ==> v1/Service
		                  NAME                   AGE
		                  busted-camel-mysql        3306/TCP  5d1h
		*/
		resource := strings.Split(services[sv], "AGE\n") //splitting after AGE
		for in, res2 := range resource {

			if in == 1 {
				moreService := strings.Split(res2, "\n")
				for _, res3 := range moreService {
					tempService := make(map[string]interface{})

					//replacing all the extras whitespace between the strings
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(res3, " ")

					//spltting the services string on whitespace
					splitservice := strings.Split(s, " ")

					tempService["name"] = splitservice[0]
					tempService["age"] = splitservice[1]

					totalServices = append(totalServices, tempService)

				}

			}
		}
	}

	return totalServices

}

func pvcStatusdetails(res *rls.GetReleaseStatusResponse, d *schema.ResourceData) []map[string]interface{} {
	var pvcs []string // getting all the pvcs available in the chart
	pvcs = make([]string, 0)
	var pvccount = 0 //for counting the no of pvcs fetched

	pattern := `==> v\d/PersistentVolumeClaim\W*\w*\W*\nNAME\s+AGE\n`

	resources := strings.Split(res.Info.Status.Resources, "\n\n") //splitting the resource string on \n\n
	for _, res := range resources {

		matched, _ := regexp.MatchString(pattern, res) //fetching string containing info of pvcs
		if matched {
			pvcs = append(pvcs, res)
			pvccount++
		}

	}
	TotalPvcCount = pvccount
	var pvcDetails []map[string]interface{} //creating a list of pvcs

	for pv := 0; pv < pvccount; pv++ {

		/*format of each received PersistentVolumeClaim as a single string
				  ==> v1/PersistentVolumeClaim
		              NAME                    AGE
		              invited-kitten-mongodb  2s

		*/

		resource := strings.Split(pvcs[pv], "AGE\n") //splitting after AGE
		for in, res2 := range resource {

			if in == 1 {
				morepvc := strings.Split(res2, "\n")
				for _, res3 := range morepvc {
					tempPvc := make(map[string]interface{})

					//replacing all the extras whitespace between the strings
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(res3, " ")

					//spltting the pods string on whitespace
					splitpvc := strings.Split(s, " ")

					tempPvc["name"] = splitpvc[0]
					tempPvc["age"] = splitpvc[1]

					pvcDetails = append(pvcDetails, tempPvc)

				}

			}
		}

	}
	return pvcDetails
}

func ingressStatusdetails(res *rls.GetReleaseStatusResponse, d *schema.ResourceData) []map[string]interface{} {
	var ingres []string // getting all the ingress available in the chart
	ingres = make([]string, 0)
	var ingrescount = 0 //for counting the no of ingress available

	pattern := `==> v\d\W*\w*\W*\d/Ingress\W*\w*\W*\nNAME\s+AGE\n`

	resources := strings.Split(res.Info.Status.Resources, "\n\n") //splitting the resource string on \n\n
	for _, res := range resources {

		matched, _ := regexp.MatchString(pattern, res) //fetching string containing info of ingress
		if matched {
			ingres = append(ingres, res)
			ingrescount++
		}

	}
	TotalIngressCount = ingrescount
	var ingressDetails []map[string]interface{}
	//creating a list of ingress

	for en := 0; en < ingrescount; en++ {

		/*format of each received ingress as a single string
				 ==> v1beta1/Ingress
		            NAME             AGE
		            example2-tomcat  2m40s
		*/

		resource := strings.Split(ingres[en], "AGE\n") //splitting after AGE
		for in, res2 := range resource {

			if in == 1 {
				moreingress := strings.Split(res2, "\n")
				for _, res3 := range moreingress {
					tempIngres := make(map[string]interface{})

					//replacing all the extras whitespace between the strings
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(res3, " ")

					//spltting the ingress string on whitespace
					splitingress := strings.Split(s, " ")

					tempIngres["name"] = splitingress[0]
					tempIngres["age"] = splitingress[1]

					ingressDetails = append(ingressDetails, tempIngres)

				}

			}
		}

	}
	return ingressDetails
}
