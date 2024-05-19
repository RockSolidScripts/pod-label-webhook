/*
Copyright 2024 ganesh.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"encoding/json"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	podlog                          = logf.Log.WithName("pod-webhook")
	_      webhook.AdmissionHandler = &PodLabel{}
)

//+kubebuilder:webhook:path=/mutate-core-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups=core,resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1

type PodLabel struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func NewPodLabel(client client.Client, scheme *runtime.Scheme) admission.Handler {
	return &PodLabel{
		Client:  client,
		Decoder: admission.NewDecoder(scheme),
	}
}

func (p *PodLabel) Handle(ctx context.Context, request admission.Request) admission.Response {

	podlog.Info("Pod create/update event")

	pod := &v1.Pod{}

	err := p.Decoder.Decode(request, pod)

	if err != nil {
		podlog.Error(err, "error decoding the admission request")
		return admission.Errored(http.StatusBadRequest, err)
	}

	podlog.Info("handling the pod CREATE/UPDATE event for", "Pod", pod.Name)

	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}

	pod.Labels["ganesh-label"] = "my-test-label"
	pod.Annotations["test"] = "test"

	marshelledPodJson, err := json.Marshal(pod)

	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	podlog.Info("the following pod has been labeled successfully", "namespace", pod.Namespace)
	return admission.PatchResponseFromRaw(request.AdmissionRequest.Object.Raw, marshelledPodJson)
}
