package main

import (
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog"
	"os"
	"time"
)

var (
	client *clientset.Clientset
)

func getNewLock(lockname, podname, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockname,
			Namespace: namespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podname,
		},
	}
}

func doStuff() {
	for {
		klog.Info("doing stuff...")
		time.Sleep(5 * time.Second)
	}
}

func runLeaderSelection(lock *resourcelock.LeaseLock, ctx context.Context, id string) {
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				doStuff()
			},
			OnStoppedLeading: func() {
				klog.Info("no longer the leader, stying inactive.")
			},
			OnNewLeader: func(current_id string) {
				if current_id == id {
					klog.Info("still the leader!")
					return
				}
				klog.Info("new leader is %s", current_id)
			},
		},
	})
}

func main() {
	var (
		leaseLockName      string
		leaseLockNamespace string
		podName            = os.Getenv("POD_NAME")
	)
	flag.StringVar(&leaseLockName, "lease-name", "", "Name of lease lock")
	flag.StringVar(&leaseLockNamespace, "lease-namespace", "default", "Name of lease lock namespace")
	flag.Parse()

	if leaseLockName == "" {
		klog.Fatal("missing lease-name flag")
	}
	if leaseLockNamespace == "" {
		klog.Fatal("missing lease-namespace flag")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal("failed to get kubeconfig")
	}
	client = clientset.NewForConfigOrDie(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lock := getNewLock(leaseLockName, podName, leaseLockNamespace)
	runLeaderSelection(lock, ctx, podName)
}