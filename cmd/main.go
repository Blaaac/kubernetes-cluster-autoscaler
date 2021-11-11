package main

import (
	"context"
	"encoding/json"
	"kubernetes-cluster-autoscaler/pkg/cloud/openstack"
	"kubernetes-cluster-autoscaler/pkg/cloud/openstack/handle-node-delete"
	"kubernetes-cluster-autoscaler/pkg/cloud/openstack/handle-node-add"
	"kubernetes-cluster-autoscaler/pkg/common/datastructures"
	"kubernetes-cluster-autoscaler/pkg/common/functions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
	"plugin"
	"sync"
)

var (
	wg                  sync.WaitGroup
	cloudType           string
	modifyEventAnalyzer func(datastructures.Event, string, string, string, string, string)
	deleteEventAnalyzer func(datastructures.Event, string, string, string, string, string)
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cloudType = openstackinit.ReadConfig()

	if cloudType != "OpenStack" {
		pluginName := cloudType + ".so"
		plugIn, err := plugin.Open("plugin/" + pluginName)
		if err != nil {
			log.Fatalf("[ERROR] Cloud not load the plugin %s in plugin directory %v", pluginName, err)
		}

		var ok bool
		modifyEventAnalyzerSymbol, err := plugIn.Lookup("ModifyEventAnalyzer")
		modifyEventAnalyzer, ok = modifyEventAnalyzerSymbol.(func(datastructures.Event, string, string, string, string, string))
		if err != nil || !ok {
			log.Fatalf("[ERROR] Something went wrong while loading plugin %v", err)
		}

		deleteEventAnalyzerSymbol, err := plugIn.Lookup("DeleteEventAnalyzer")
		deleteEventAnalyzer, ok = deleteEventAnalyzerSymbol.(func(datastructures.Event, string, string, string, string, string))
		if err != nil || !ok {
			log.Fatalf("[ERROR] Something went wrong while loading plugin %v", err)
		}
		log.Printf("[INFO] %s plugin loaded succesfully", cloudType)
	}

	config := functions.LoadKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	checkErr(err)

	resource := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	w, err := dynamicClient.Resource(resource).Namespace("").Watch(context.TODO(), metav1.ListOptions{})
	checkErr(err)

	defer w.Stop()
	wg.Add(1)
	log.Println("[INFO] K8s cluster auto scalar started")
  handlenodeadd.TriggerAddNode("m1.large")
	for {
		wCh := w.ResultChan()

		for event := range wCh {
			eventFilter(event, config)
		}
	}
}

func eventFilter(event watch.Event, config *rest.Config) {
	defer handlenodeadd.PanicRecovery()

	b, _ := json.Marshal(event)
	var EventList datastructures.Event
	err := json.Unmarshal(b, &EventList)
	if err != nil {
		panic("Need to write the code. Stop containing")
	}

	switch EventList.Type {
	case "MODIFIED":
		if cloudType == "OpenStack" {
			handlenodeadd.ModifyEventAnalyzer(EventList, config)
		} else {
			modifyEventAnalyzer(EventList, openstackinit.ProjectName, openstackinit.ClientSecret, openstackinit.ClientID, openstackinit.AWSRegion, openstackinit.AuthFile)
		}
	case "DELETED":
		if cloudType == "OpenStack" {
			handlenodedelete.DeleteEventAnalyzer(EventList, config)
		} else {
			deleteEventAnalyzer(EventList, openstackinit.ProjectName, openstackinit.ClientSecret, openstackinit.ClientID, openstackinit.AWSRegion, openstackinit.AuthFile)
		}
	}
}
